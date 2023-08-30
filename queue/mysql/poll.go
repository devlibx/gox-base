package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/oklog/ulid/v2"
	errors1 "github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

var useV1 = true

func (t *topRowFinder) Start() {
	if useV1 {
		return
	}
	t.internalStart(true)
	go t.internalStart(false)
}

func (t *topRowFinder) Stop() {
	t.stop = true
}

func (t *topRowFinder) internalStart(runOnce bool) {
	errorCount := 1

	query := "select MIN(id) FROM jobs WHERE tenant=? AND state=? AND job_type=?"
	query = t.queryRewriter.RewriteQuery("jobs", query)
	if t.usePreparedStatement && t.findTopProcessAtQueryStmt == nil {
		t.findTopProcessAtQueryStmt, _ = t.db.PrepareContext(context.Background(), query)
	}

	for {
		var _time string
		var err error
		ctxQ, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		t.logger.Debug("[start] find the top process at from queue", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
		if t.findTopProcessAtQueryStmt != nil {
			err = t.findTopProcessAtQueryStmt.QueryRowContext(ctxQ, t.tenant, queue.StatusScheduled, t.jobType).Scan(&_time)
		} else {
			err = t.db.QueryRowContext(ctxQ, query, t.tenant, queue.StatusScheduled, t.jobType).Scan(&_time)
		}
		cancelFunc()

		t.logger.Debug("[end] find the top process at from queue", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
		if err == nil && _time != "" {
			a := ulid.MustParse(_time)
			t.smallestProcessAtTime = time.UnixMilli(int64(a.Time()))
			// t.smallestProcessAtTime, _ = time.Parse("2006-01-02 15:04:05", _time)

			// We are done return - this has to run only once at boot
			// NOTE - runOnce is true only at the time of boot
			if runOnce {
				return
			}

			// Wait for sometime
		noOp:
			select {
			case <-time.After(1 * time.Second):
				break // No op
			case ev, closed := <-t.refreshChannel:
				if closed {
					if ev.time.IsZero() {
						// time.Sleep(10 * time.Millisecond)
					} else {
						if t.smallestProcessAtTime.After(ev.time) {
							// time.Sleep(10 * time.Millisecond)
							t.logger.Debug("* ignore - request to find the top process at from queue *")
							goto noOp
						} else {
							t.logger.Debug("request to find the top process at from queue", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant), zap.String("time", ev.time.Local().String()))
						}
					}
				} else {
					time.Sleep(10 * time.Millisecond)
				}
			}

		} else {
			t.logger.Debug("not able to fetch the smallest value", zap.Error(err))
			time.Sleep(1 * time.Second)
			errorCount++
			if t.smallestProcessAtTime.IsZero() && errorCount%9 == 0 {
				t.logger.Info("[WARN - expected at boot-up] did not find smallest processed at time or if there is not job for given job type", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
			}
		}

		if t.stop {
			t.logger.Info("Stopping top row finder", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
			break
		}

		if runOnce {
			break
		}
	}
}

func (q *queueImpl) Poll(ctx context.Context, req queue.PollRequest) (*queue.PollResponse, error) {
	if useV1 {
		return q.internalPollV1(ctx, req)
	}
	return q.internalPoll(ctx, req)
}

func (q *queueImpl) internalPoll(ctx context.Context, req queue.PollRequest) (result *queue.PollResponse, err error) {
	var rowFinder *topRowFinder
	var ok bool
	if rowFinder, ok = q.topRowFinderCron[req.JobType]; !ok {
		return nil, errors.New("failed to find top row for jobType=%d tenant=%d", req.JobType, req.Tenant)
	} else if rowFinder.smallestProcessAtTime.IsZero() {
		return nil, errors.New("[WARN - expected at boot-up time] failed to find top row for jobType=%d tenant=%d", req.JobType, req.Tenant)
	}

	// Begin a transaction
	tx, err := q.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin txn to schedule job")
	}

	// Rollback or commit at the end
	defer func() {
		if p := recover(); p != nil {
			q.logger.Error("found error in polling", zap.Any("error", p))
			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	processAt := rowFinder.smallestProcessAtTime.Truncate(time.Second)
	remainingRetries := 0
	partitionTime := partitionBasedOnProcessAtTime(processAt)
	result = &queue.PollResponse{RecordPartitionTime: partitionTime, ProcessAtTimeUsed: processAt}

	pollQuery := "SELECT id, pending_execution FROM jobs WHERE process_at=? AND tenant=? AND job_type=? AND state=? AND part=? LIMIT 1 FOR UPDATE SKIP LOCKED"
	pollQuery = q.queryRewriter.RewriteQuery("jobs", pollQuery)

	updatePollResultQuery := "UPDATE jobs SET state=?, version=version+1, pending_execution=pending_execution-1 WHERE id=? AND part=?"
	updatePollResultQuery = q.queryRewriter.RewriteQuery("jobs", updatePollResultQuery)

	if q.usePreparedStatement {
		var failed string
		q.pollQueryStatementInitOnce.Do(func() {
			if q.pollQueryStatement, err = q.db.Prepare(pollQuery); err != nil {
				failed = "poll query statement"
			} else if q.updatePollRecordStatement, err = q.db.Prepare(updatePollResultQuery); err != nil {
				failed = "poll update result statement"
			}
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create %s", failed)
		}
		err = tx.StmtContext(ctx, q.pollQueryStatement).QueryRowContext(ctx, processAt, req.Tenant, req.JobType, queue.StatusScheduled, partitionTime).Scan(&result.Id, &remainingRetries)
	} else {
		err = tx.QueryRowContext(ctx, pollQuery, processAt, req.Tenant, req.JobType, queue.StatusScheduled, partitionTime).Scan(&result.Id, &remainingRetries)
	}
	if errors1.Is(err, sql.ErrNoRows) {
		err = errors.Wrap(queue.NoJobsToRunAtCurrently, "jobType=%d tenant=%d topTime=%s", req.JobType, req.Tenant, processAt.String())

		// Push event to make sure we refresh asap
		if rf, ok := q.topRowFinderCron[req.JobType]; ok {
			rf.refreshChannel <- refreshEvent{time: processAt}
		}

		return nil, err
	} else if err != nil {
		err = fmt.Errorf("failed to find top row for jobType=%d tenant=%d topTime=%s err=%w", req.JobType, req.Tenant, processAt.String(), err)
		return nil, err
	}

	if remainingRetries <= 0 {
		if _, err = tx.ExecContext(ctx, "UPDATE jobs SET state=?, sub_state=? WHERE id=?", queue.StatusFailed, queue.SubStatusNoRetryPendingError, result.Id); err != nil {
			return nil, fmt.Errorf("we found a record which is already dead: id=%s mark it failed and set it retries < 0", result.Id)
		}
	}

	// Update the row within the same transaction
	var res sql.Result
	if q.usePreparedStatement {
		res, err = tx.StmtContext(ctx, q.updatePollRecordStatement).ExecContext(ctx, queue.StatusProcessing, result.Id, partitionTime)
	} else {
		// res, err = tx.ExecContext(ctx, "UPDATE jobs SET state=?, version=version+1, pending_execution=pending_execution-1 WHERE id=? AND part=?", queue.StatusProcessing, result.Id, partitionTime)
		res, err = tx.ExecContext(ctx, updatePollResultQuery, queue.StatusProcessing, result.Id, partitionTime)
	}
	var t int64
	if err != nil {
		err = fmt.Errorf("failed to update the job table pending_execution: %w id=%s", err, result.Id)
		return nil, err
	} else if t, err = res.RowsAffected(); err == nil && t == 0 {
		err = fmt.Errorf("failed to update the job table (concurrent update): %w id=%s", err, result.Id)
		return nil, err
	}

	return
}
