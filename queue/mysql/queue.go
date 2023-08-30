package queue

import (
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"go.uber.org/zap"
	"sync"
	"time"
)

type queueImpl struct {
	cf gox.CrossFunction

	db *sql.DB

	smallestProcessedAt map[string]time.Time

	pollQueryStatement         *sql.Stmt
	updatePollRecordStatement  *sql.Stmt
	pollQueryStatementInitOnce *sync.Once

	initOnce  *sync.Once
	closeOnce *sync.Once

	storeBackend  queue.StoreBackend
	queueConfig   queue.MySqlBackedQueueConfig
	idGenerator   queue.IdGenerator
	queryRewriter queue.QueryRewriter

	logger *zap.Logger

	topRowFinderCronMutex *sync.RWMutex
	topRowFinderCron      map[int]*topRowFinder

	jobTypeRowInfo map[int]*jobTypeRowInfo

	usePreparedStatement bool
}

type refreshEvent struct {
	time time.Time
}

type topRowFinder struct {
	jobType               int
	tenant                int
	smallestProcessAtTime time.Time
	db                    *sql.DB
	logger                *zap.Logger
	stop                  bool
	queryRewriter         queue.QueryRewriter
	refreshChannel        chan refreshEvent

	usePreparedStatement          bool
	findTopProcessAtQueryStmt     *sql.Stmt
	findTopProcessAtQueryStmtOnce *sync.Once
}

func NewQueue(cf gox.CrossFunction, storeBackend queue.StoreBackend, queueConfig queue.MySqlBackedQueueConfig, idGenerator queue.IdGenerator, queryRewriter queue.QueryRewriter) (*queueImpl, error) {

	// Get a DB to be used
	db, err := storeBackend.GetSqlDb()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build new mysql backed queue. Could not get sql.Db from store backend")
	}

	q := &queueImpl{
		db:            db,
		storeBackend:  storeBackend,
		queueConfig:   queueConfig,
		idGenerator:   idGenerator,
		queryRewriter: queryRewriter,
		logger:        cf.Logger().Named("scheduler"),

		topRowFinderCron:      map[int]*topRowFinder{},
		topRowFinderCronMutex: &sync.RWMutex{},

		pollQueryStatementInitOnce: &sync.Once{},

		usePreparedStatement: queueConfig.UsePreparedStatement,

		jobTypeRowInfo: map[int]*jobTypeRowInfo{},
	}

	// Run job top finder - we can configure max job type id
	if queueConfig.MaxJobType <= 0 {
		queueConfig.MaxJobType = 1
	}
	for i := 1; i <= queueConfig.MaxJobType && !queueConfig.DontRunPoller; i++ {
		q.topRowFinderCron[i] = &topRowFinder{
			db:                            db,
			jobType:                       i,
			tenant:                        queueConfig.Tenant,
			logger:                        q.logger,
			queryRewriter:                 queryRewriter,
			usePreparedStatement:          q.usePreparedStatement,
			refreshChannel:                make(chan refreshEvent, 2),
			findTopProcessAtQueryStmtOnce: &sync.Once{},
		}
		q.topRowFinderCron[i].refreshChannel <- refreshEvent{}
		q.topRowFinderCron[i].Start()
	}

	for i := 1; i <= queueConfig.MaxJobType && !queueConfig.DontRunPoller; i++ {
		q.jobTypeRowInfo[i] = &jobTypeRowInfo{
			jobType:                         i,
			tenant:                          queueConfig.Tenant,
			db:                              db,
			logger:                          cf.Logger().Named(fmt.Sprintf("scheduler-jobType=%d-tenant=%d", i, queueConfig.Tenant)),
			usePreparedStatement:            q.usePreparedStatement,
			queryRewriter:                   queryRewriter,
			smallestScheduledJobProcessTime: time.Time{},
			findSmallestScheduledJobProcessAtTimeStatementOnce: &sync.Once{},
			smallestScheduledJobProcessTimeLock:                &sync.RWMutex{},
		}
		if err = q.jobTypeRowInfo[i].Init(); err != nil {
			return nil, err
		}
	}

	return q, nil
}

func partitionBasedOnProcessAtTime(t time.Time) time.Time {
	return endOfWeek(t)
}

func endOfWeekPlusOneWeek(inputTime time.Time) time.Time {
	return endOfWeek(inputTime).Add(7 * 24 * time.Hour)
}

func endOfWeek(inputTime time.Time) time.Time {
	inputTime = inputTime.Truncate(time.Hour).Add(time.Duration(-1 * inputTime.Hour()))
	daysUntilSunday := 0
	switch inputTime.Weekday() {
	case time.Monday:
		daysUntilSunday = 6
		break
	case time.Tuesday:
		daysUntilSunday = 5
		break
	case time.Wednesday:
		daysUntilSunday = 4
		break
	case time.Thursday:
		daysUntilSunday = 3
		break
	case time.Friday:
		daysUntilSunday = 2
		break
	case time.Saturday:
		daysUntilSunday = 1
		break
	case time.Sunday:
		daysUntilSunday = 0
		break
	}

	// Use the Add method to add the remaining days to the input time.

	endOfWeekTime := inputTime.Add(time.Duration(daysUntilSunday) * 24 * time.Hour)

	// Set the time to the end of the day (23:59:59).
	endOfWeekTime = time.Date(endOfWeekTime.Year(), endOfWeekTime.Month(), endOfWeekTime.Day(), 23, 59, 59, 0, time.Now().Location())

	return endOfWeekTime
}
