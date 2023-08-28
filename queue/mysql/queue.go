package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/sethvargo/go-retry"
	"sync"
	"time"
)

type queueImpl struct {
	db *sql.DB

	insertStatement *sql.Stmt

	smallestProcessedAt time.Time

	initOnce     *sync.Once
	closeOnce    *sync.Once
	storeBackend queue.StoreBackend
	queueConfig  queue.MySqlBackedQueueConfig

	idGenerator   queue.IdGenerator
	queryRewriter queue.QueryRewriter
}

func (q *queueImpl) Schedule(ctx context.Context, req queue.ScheduleRequest) (result *queue.ScheduleResponse, err error) {
	if err = retry.Exponential(ctx, 1*time.Second, func(ctx context.Context) error {
		if result, err = q.internalSchedule(ctx, req); err != nil {
			fmt.Println("TODO - put retry")
			return err
		}
		return err
	}); err != nil {
		err = errors.Wrap(err, "failed to schedule to mysql queue: %v", req)
	}
	return
}

func (q *queueImpl) internalSchedule(ctx context.Context, req queue.ScheduleRequest) (*queue.ScheduleResponse, error) {
	processAt := req.At.Truncate(time.Second)
	id := q.idGenerator.GenerateId(processAt)

	// Min count = 1 i.e. each row is processed min once
	remainingExecution := req.RemainingExecution
	if remainingExecution <= 0 {
		remainingExecution = 1
	}

	state := queue.StatusScheduled
	subState := queue.SubStatusScheduledOk

	// Metadata
	var properties interface{}
	properties = `{"": ""}`
	if req.Properties != nil {
		properties = req.Properties
		if jsonData, err := json.Marshal(req.Properties); err == nil {
			properties = jsonData
		} else {
			return nil, fmt.Errorf("failed to persist (metadata is bad): %w", err)
		}
	}

	// id, correlation_id, job_type, process_at, state, sub_state, properties, version, pending_execution, string_udf_1, string_udf_2, int_udf_1, int_udf_2
	// Execute the INSERT statement
	if _, err := q.insertStatement.Exec(
		id, req.CorrelationId, req.JobType, processAt, state, subState, properties, 1, remainingExecution, req.StringUdf1, req.StringUdf2, req.IntUdf1, req.IntUdf1,
	); err != nil {
		return nil, errors.Wrap(err, "failed to schedule: %v", req)
	} else {
		return &queue.ScheduleResponse{Id: id}, err
	}
}

func (s *queueImpl) topRowScanner(jobType string) {
	for {
		var t string
		err := s.db.QueryRow("SELECT process_at FROM jobs WHERE state=? AND job_type=? order by process_at asc", queue.StatusScheduled, jobType).Scan(&t)
		if err == nil {
			s.smallestProcessedAt, err = time.Parse("2006-01-02 15:04:05", t)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				fmt.Println("[WARN] failed to find smallest processed at time", err, t)
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			time.Sleep(1 * time.Second)
			fmt.Println("[WARN - expected at boot-up] did not find smallest processed at time", err, t)
		}
	}
}

func NewQueue(storeBackend queue.StoreBackend, queueConfig queue.MySqlBackedQueueConfig, idGenerator queue.IdGenerator, queryRewriter queue.QueryRewriter) (*queueImpl, error) {

	// Get a DB to be used
	db, err := storeBackend.GetSqlDb()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build new mysql backed queue. Could not get sql.Db from store backend")
	}

	// Generate insert statment
	insertQuery := `
			INSERT INTO jobs (
				id, correlation_id, job_type, process_at, state, sub_state, properties, version, pending_execution, string_udf_1, string_udf_2, int_udf_1, int_udf_2			                       
			) 
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)			
			`
	insertQuery = queryRewriter.RewriteQuery(insertQuery)
	insertStmt, err := db.Prepare(insertQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build insert schedule query statement")
	}

	q := &queueImpl{
		db:              db,
		insertStatement: insertStmt,
		storeBackend:    storeBackend,
		queueConfig:     queueConfig,
		idGenerator:     idGenerator,
		queryRewriter:   queryRewriter,
	}

	// Populate to the smallest row
	go q.topRowScanner("cron")

	return q, nil
}
