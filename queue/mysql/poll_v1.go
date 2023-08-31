package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/oklog/ulid/v2"
	pkgErrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

type jobTypeRowInfo struct {
	jobType int
	tenant  int

	usePreparedStatement bool
	queryRewriter        queue.QueryRewriter
	db                   *sql.DB
	logger               *zap.Logger

	smallestScheduledJobProcessTimeLock *sync.RWMutex

	smallestScheduledJobProcessTime                    time.Time
	findSmallestScheduledJobProcessAtTimeStatement     *sql.Stmt
	findSmallestScheduledJobProcessAtTimeStatementOnce *sync.Once
}

func (j *jobTypeRowInfo) Init() (err error) {
	j.findSmallestScheduledJobProcessAtTimeStatementOnce.Do(func() {
		query := "select MIN(id) FROM jobs WHERE tenant=? AND state=? AND job_type=?"
		query = j.queryRewriter.RewriteQuery("jobs", query)
		if j.findSmallestScheduledJobProcessAtTimeStatement, err = j.db.PrepareContext(context.Background(), query); err != nil {
			err = errors.Wrap(err, "failed to prepare smallest scheduled job processAt time statement")
		}
	})
	return
}

func (j *jobTypeRowInfo) ensureSmallestScheduledJobProcessTime(ctx context.Context) (err error) {
	var id sql.NullString
	var i ulid.ULID
	ctxQ, cancelFunc := context.WithTimeout(ctx, 1*time.Second)
	defer cancelFunc()
	err = j.findSmallestScheduledJobProcessAtTimeStatement.QueryRowContext(ctxQ, j.tenant, queue.StatusScheduled, j.jobType).Scan(&id)
	if err == nil && id.Valid {
		if i, err = ulid.Parse(id.String); err == nil {
			j.smallestScheduledJobProcessTime = time.UnixMilli(int64(i.Time()))
		} else {
			err = errors.Wrap(err, "failed to get time from id: id=%s", i)
		}
	} else {
		err = errors.Wrap(err, "failed to find the smallest scheduled job processAt time: tenant=%d, jobType=%d, state=%d", j.tenant, j.jobType, queue.StatusScheduled)
	}
	return
}

func (q *queueImpl) initPollQueriesV1(ctx context.Context) (pollQuery string, updatePollResultQuery string, err error) {

	// Build poll query with table rewrite
	if q.useMinQueryToPickLatestRow {
		pollQuery = "select MIN(id) from jobs where tenant=? AND state=? AND job_type=?  for update SKIP LOCKED"
	} else {
		pollQuery = "SELECT id, pending_execution FROM jobs WHERE process_at=? AND tenant=? AND job_type=? AND state=? AND part=? LIMIT 1 FOR UPDATE SKIP LOCKED"
	}
	pollQuery = q.queryRewriter.RewriteQuery("jobs", pollQuery)

	// Build update query with table rewrite
	updatePollResultQuery = "UPDATE jobs SET state=?, version=version+1, pending_execution=pending_execution-1 WHERE id=? AND part=?"
	updatePollResultQuery = q.queryRewriter.RewriteQuery("jobs", updatePollResultQuery)

	if q.usePreparedStatement {
		var failed string
		q.pollQueryStatementInitOnce.Do(func() {
			if q.pollQueryStatement, err = q.db.PrepareContext(ctx, pollQuery); err != nil {
				failed = "poll query statement"
			} else if q.updatePollRecordStatement, err = q.db.PrepareContext(ctx, updatePollResultQuery); err != nil {
				failed = "poll update result statement"
			}
		})
		if err != nil {
			return pollQuery, updatePollResultQuery, errors.Wrap(err, "failed to create %s", failed)
		}
	}
	return
}

func (q *queueImpl) ensureWeMarkJobsWithNoRemainingTryProperly(ctx context.Context, tx *sql.Tx, id string) (err error) {
	if _, err = tx.ExecContext(ctx, "UPDATE jobs SET state=?, sub_state=? WHERE id=?", queue.StatusFailed, queue.SubStatusNoRetryPendingError, id); err != nil {
		return fmt.Errorf("we found a record which is already dead: id=%s mark it failed and set it retries < 0", id)
	}
	return
}

func (q *queueImpl) internalPollV1(ctx context.Context, req queue.PollRequest) (result *queue.PollResponse, err error) {

	// Step 1 - make sure we have configured this job type - jobTypeRowInfo contains the smallest time for this job type and tenant
	var jobTypeRowInfoObj *jobTypeRowInfo
	var ok bool
	if jobTypeRowInfoObj, ok = q.jobTypeRowInfo[req.JobType]; !ok {
		return nil, errors.New("the job type given in request is not valid - check your MySqlBackedQueueConfig.max_job_type config value")
	}

	// Make sure we have poll and update query
	var pollQuery, updatePollResultQuery string
	pollQuery, updatePollResultQuery, err = q.initPollQueriesV1(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build poll and update query")
	}

	// Begin a transaction
	tx, err := q.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin txn to schedule job")
	}
	defer func() {
		if p := recover(); p != nil {
			q.logger.Error("found error in polling", zap.Any("error", p))
			if e := tx.Rollback(); e != nil {
				q.logger.Error("something is wrong - tx failed to rollback after panic")
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				q.logger.Error("something is wrong - tx failed to rollback")
			}
		} else {
			if e := tx.Commit(); e != nil {
				q.logger.Error("something is wrong - tx failed to commit")
			}
		}
	}()

	// Ensure we have the smallest scheduled time - read lock to see if we have it already
	processAt := time.Time{}
	jobTypeRowInfoObj.smallestScheduledJobProcessTimeLock.RLock()
	if !jobTypeRowInfoObj.smallestScheduledJobProcessTime.IsZero() {
		processAt = jobTypeRowInfoObj.smallestScheduledJobProcessTime.Truncate(time.Second)
	}
	jobTypeRowInfoObj.smallestScheduledJobProcessTimeLock.RUnlock()

	// If it is missing then force refresh
	tryCount := 0
tryRefreshingSmallestScheduledJobProcessTime:
	if processAt.IsZero() {
		jobTypeRowInfoObj.smallestScheduledJobProcessTimeLock.Lock()
		if jobTypeRowInfoObj.smallestScheduledJobProcessTime.IsZero() {
			if err = jobTypeRowInfoObj.ensureSmallestScheduledJobProcessTime(ctx); err == nil {
				processAt = jobTypeRowInfoObj.smallestScheduledJobProcessTime.Truncate(time.Second)
			}
		} else if processAt.IsZero() && !jobTypeRowInfoObj.smallestScheduledJobProcessTime.IsZero() {
			if err = jobTypeRowInfoObj.ensureSmallestScheduledJobProcessTime(ctx); err == nil {
				processAt = jobTypeRowInfoObj.smallestScheduledJobProcessTime.Truncate(time.Second)
			}
		} else {
			processAt = jobTypeRowInfoObj.smallestScheduledJobProcessTime.Truncate(time.Second)
		}
		jobTypeRowInfoObj.smallestScheduledJobProcessTimeLock.Unlock()
	}
	if err != nil {
		if tryCount == 1 {
			err = errors.Wrap(err, "[retry] failed to find the smallest scheduled processAt time")
		} else {
			err = errors.Wrap(err, "failed to find the smallest scheduled processAt time")
		}
		return
	}

	// Params to query for top job
	remainingRetries := 0
	partitionTime := partitionBasedOnProcessAtTime(processAt)
	result = &queue.PollResponse{RecordPartitionTime: partitionTime, ProcessAtTimeUsed: processAt}

	if q.usePreparedStatement && q.useMinQueryToPickLatestRow {
		var resultId sql.NullString
		remainingRetries = 1
		err = tx.StmtContext(ctx, q.pollQueryStatement).
			QueryRowContext(ctx, req.Tenant, req.JobType, queue.StatusScheduled).
			Scan(&resultId)
		if err == nil && resultId.Valid {
			result.Id = resultId.String
		} else {
			err = errors.Wrap(sql.ErrNoRows, "no rows found with min(id)")
		}
	} else if q.usePreparedStatement {
		err = tx.StmtContext(ctx, q.pollQueryStatement).
			QueryRowContext(ctx, processAt, req.Tenant, req.JobType, queue.StatusScheduled, partitionTime).
			Scan(&result.Id, &remainingRetries)
	} else {
		err = tx.QueryRowContext(ctx, pollQuery, processAt, req.Tenant, req.JobType, queue.StatusScheduled, partitionTime).
			Scan(&result.Id, &remainingRetries)
	}

	if pkgErrors.Is(err, sql.ErrNoRows) {
		// Once try to get the smallest process time if we got no rows
		if tryCount == 0 {
			tryCount = 1
			processAt = time.Time{}
			goto tryRefreshingSmallestScheduledJobProcessTime
		} else {
			err = errors.Wrap(queue.NoJobsToRunAtCurrently, "jobType=%d tenant=%d topTime=%s", req.JobType, req.Tenant, processAt.String())
			return
		}
	} else if err != nil {
		err = fmt.Errorf("failed to find top row for jobType=%d tenant=%d topTime=%s err=%w", req.JobType, req.Tenant, processAt.String(), err)
		return
	}

	// Just in case we have records which are still in scheduled with zero retries - fix them
	if remainingRetries <= 0 {
		if err = q.ensureWeMarkJobsWithNoRemainingTryProperly(ctx, tx, result.Id); err != nil {
			return
		}
	}

	// Update the row within the same transaction
	var res sql.Result
	if q.usePreparedStatement {
		res, err = tx.StmtContext(ctx, q.updatePollRecordStatement).ExecContext(ctx, queue.StatusProcessing, result.Id, partitionTime)
	} else {
		res, err = tx.ExecContext(ctx, updatePollResultQuery, queue.StatusProcessing, result.Id, partitionTime)
	}

	var noOfUpdatedRecords int64
	if err != nil {
		err = fmt.Errorf("failed to update the job table pending_execution: %w id=%s", err, result.Id)
	} else if noOfUpdatedRecords, err = res.RowsAffected(); err == nil && noOfUpdatedRecords == 0 {
		err = fmt.Errorf("failed to update the job table (concurrent update): %w id=%s", err, result.Id)
	}

	return
}

func (q *queueImpl) idToTime(id string) (result time.Time, err error) {
	var i ulid.ULID
	if i, err = ulid.Parse(id); err == nil {
		result = time.UnixMilli(int64(i.Time()))
	} else {
		err = errors.Wrap(err, "failed to get time from id: id=%s", i)
	}
	return result, err
}
