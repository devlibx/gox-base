package queue

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/bombsimon/mysql-error-numbers"
	mysqlerrnum "github.com/bombsimon/mysql-error-numbers"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/devlibx/gox-base/serialization"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
	"time"
)

func (q *queueImpl) Schedule(ctx context.Context, req queue.ScheduleRequest) (result *queue.ScheduleResponse, err error) {
	if err = retry.Exponential(ctx, 1*time.Second, func(ctx context.Context) error {
		if result, err = q.internalScheduleV1(ctx, req, req.InternalTx); err != nil {
			var e *mysql.MySQLError
			if errors.As(err, &e) && e.Number == mysqlerrnum.ER_LOCK_WAIT_TIMEOUT {
				q.logger.Info("[retry] error in scheduling job", zap.String("error", e.Error()))
				return nil
			}
			return err
		}
		return err
	}); err != nil {
		err = errors.Wrap(err, "failed to schedule to mysql queue: %v", req)
	}
	return
}

func (q *queueImpl) internalScheduleV1(ctx context.Context, req queue.ScheduleRequest, tx *sql.Tx) (result *queue.ScheduleResponse, err error) {
	processAt := req.At.Truncate(time.Second)
	id := q.idGenerator.GenerateId(processAt)

	// Min count = 1 i.e. each row is processed min once
	remainingExecution := req.RemainingExecution
	if remainingExecution <= 0 {
		remainingExecution = 1
	}

	// Initial state and sub-state
	state := queue.StatusScheduled
	subState := queue.SubStatusScheduledOk

	// We get the partition based on the process At - by default it is end of next week
	archiveAfter := queue.InternalImplEndOfWeek(processAt)

	// Metadata
	properties := `{"": ""}`
	if req.Properties != nil {
		// properties = req.Properties
		if properties, err = serialization.Stringify(req.Properties); err != nil {
			return nil, fmt.Errorf("failed to persist (metadata is bad): %w", err)
		}
	}

	// Generate insert job data statement
	insertJobDataQuery := `
			INSERT INTO jobs_data
				(id, tenant, string_udf_1, string_udf_2, int_udf_1, int_udf_2, properties, retry_group, part)
			VALUES 
			    (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	insertJobDataQuery = q.queryRewriter.RewriteQuery("jobs_data", insertJobDataQuery)

	insertJobQuery := `
			INSERT INTO jobs 
			    (id, tenant, correlation_id, job_type, process_at, state, sub_state, version, pending_execution, part) 
			VALUES
			    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)			
	`
	insertJobQuery = q.queryRewriter.RewriteQuery("jobs", insertJobQuery)

	q.insertJobStatementOnce.Do(func() {
		if q.insertJobStatement, err = q.db.PrepareContext(ctx, insertJobQuery); err != nil {
			err = errors.Wrap(err, "Failed to build insert job query")
		} else if q.insertJobDataStatement, err = q.db.PrepareContext(ctx, insertJobDataQuery); err != nil {
			err = errors.Wrap(err, "Failed to build insert job data query")
		}
	})
	if err != nil {
		return nil, err
	}

	if tx != nil {
		if req.InternalRetryGroupId == "" {
			req.InternalRetryGroupId = uuid.NewString()
		}
		if _, err = tx.StmtContext(ctx, q.insertJobDataStatement).ExecContext(ctx, id, req.Tenant, req.StringUdf1, req.StringUdf2, req.IntUdf1, req.IntUdf2, properties, req.InternalRetryGroupId, archiveAfter); err != nil {
			return nil, errors.Wrap(err, "failed to schedule (insert job data failed): %v", req)
		}
	} else {
		if req.InternalRetryGroupId == "" {
			req.InternalRetryGroupId = uuid.NewString()
		}
		if _, err = q.insertJobDataStatement.ExecContext(ctx, id, req.Tenant, req.StringUdf1, req.StringUdf2, req.IntUdf1, req.IntUdf2, properties, req.InternalRetryGroupId, archiveAfter); err != nil {
			// if _, err = q.db.ExecContext(ctx, insertJobDataQuery, id, req.Tenant, req.StringUdf1, req.StringUdf2, req.IntUdf1, req.IntUdf2, properties, archiveAfter); err != nil {
			return nil, errors.Wrap(err, "failed to schedule (insert job data failed): %v", req)
		}
	}

	if tx != nil {
		if _, err = tx.StmtContext(ctx, q.insertJobStatement).ExecContext(ctx, id, req.Tenant, req.CorrelationId, req.JobType, processAt, state, subState, 1, remainingExecution, archiveAfter); err != nil {
			return nil, errors.Wrap(err, "failed to schedule (insert job failed): %v", req)
		}
	} else {
		if _, err = q.insertJobStatement.ExecContext(ctx, id, req.Tenant, req.CorrelationId, req.JobType, processAt, state, subState, 1, remainingExecution, archiveAfter); err != nil {
			// if _, err = q.db.ExecContext(ctx, insertJobQuery, id, req.Tenant, req.CorrelationId, req.JobType, processAt, state, subState, 1, remainingExecution, archiveAfter); err != nil {
			return nil, errors.Wrap(err, "failed to schedule (insert job failed): %v", req)
		}
	}

	result = &queue.ScheduleResponse{Id: id}
	return
}
