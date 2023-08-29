package queue

import (
	"database/sql"
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

	initOnce  *sync.Once
	closeOnce *sync.Once

	storeBackend  queue.StoreBackend
	queueConfig   queue.MySqlBackedQueueConfig
	idGenerator   queue.IdGenerator
	queryRewriter queue.QueryRewriter

	logger *zap.Logger

	topRowFinderCronMutex *sync.RWMutex
	topRowFinderCron      map[int]*topRowFinder
}

type topRowFinder struct {
	jobType               int
	tenant                int
	smallestProcessAtTime time.Time
	db                    *sql.DB
	logger                *zap.Logger
	stop                  bool
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
	}

	// Run job top finder - we can configure max job type id
	if queueConfig.MaxJobType <= 0 {
		queueConfig.MaxJobType = 1
	}
	for i := 1; i <= queueConfig.MaxJobType; i++ {
		q.topRowFinderCron[i] = &topRowFinder{db: db, jobType: i, tenant: queueConfig.Tenant, logger: q.logger}
		q.topRowFinderCron[i].Start()
	}

	return q, nil
}
