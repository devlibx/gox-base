package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/devlibx/gox-base/v2/queue"
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

	insertJobStatement     *sql.Stmt
	insertJobDataStatement *sql.Stmt
	insertJobStatementOnce *sync.Once

	initOnce  *sync.Once
	closeOnce *sync.Once

	storeBackend  queue.StoreBackend
	queueConfig   queue.MySqlBackedQueueConfig
	idGenerator   queue.IdGenerator
	queryRewriter queue.QueryRewriter

	logger *zap.Logger

	jobTypeRowInfo map[int]*jobTypeRowInfo

	readJobDetailsOnce          *sync.Once
	readJobDetailsStatement     *sql.Stmt
	readJobDataDetailsStatement *sql.Stmt
	updateJobStatusStatement    *sql.Stmt
	updateJobDataStatement      *sql.Stmt

	usePreparedStatement       bool
	useMinQueryToPickLatestRow bool
}

type refreshEvent struct {
	time time.Time
}

func NewQueue(cf gox.CrossFunction, storeBackend queue.StoreBackend, queueConfig queue.MySqlBackedQueueConfig, idGenerator queue.IdGenerator, queryRewriter queue.QueryRewriter) (*queueImpl, error) {

	// Get a DB to be used
	db, err := storeBackend.GetSqlDb()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build new mysql backed queue. Could not get sql.Db from store backend")
	}

	// Setup ID generator (if not provided) - this will generate ULID based IDs
	if idGenerator == nil {
		if idGenerator, err = queue.NewTimeBasedIdGenerator(); err != nil {
			err = errors.Wrap(err, "id generator is not provided and not able to create default id generator")
			return nil, err
		}
	}

	q := &queueImpl{
		db:            db,
		storeBackend:  storeBackend,
		queueConfig:   queueConfig,
		idGenerator:   idGenerator,
		queryRewriter: queryRewriter,
		logger:        cf.Logger().Named("scheduler"),

		pollQueryStatementInitOnce: &sync.Once{},

		usePreparedStatement:       queueConfig.UsePreparedStatement,
		useMinQueryToPickLatestRow: queueConfig.UseMinQueryToPickLatestRow,

		jobTypeRowInfo: map[int]*jobTypeRowInfo{},

		readJobDetailsOnce:     &sync.Once{},
		insertJobStatementOnce: &sync.Once{},
	}

	// Run job top finder - we can configure max job type id
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

	if _, _, err = q.initPollQueriesV1(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed to init poller queries at time of queue creation")
	}

	if err = q.jobInfoInit(); err != nil {
		return nil, errors.Wrap(err, "failed to init job info queries at time of queue creation")
	}

	return q, nil
}
