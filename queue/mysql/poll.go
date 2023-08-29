package queue

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"go.uber.org/zap"
	"time"
)

func (t *topRowFinder) Start() {
	go t.internalStart()
}

func (t *topRowFinder) Stop() {
	t.stop = true
}

func (t *topRowFinder) internalStart() {
	for {
		var _time string
		err := t.db.QueryRow("SELECT process_at FROM jobs WHERE tenant=? AND state=? AND job_type=? order by process_at asc", t.tenant, queue.StatusScheduled, t.jobType).Scan(&_time)
		if err == nil && _time != "" {
			t.smallestProcessAtTime, _ = time.Parse("2006-01-02 15:04:05", _time)
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
			if t.smallestProcessAtTime.IsZero() {
				t.logger.Info("[WARN - expected at boot-up] did not find smallest processed at time or if there is not job for given job type", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
			}
		}

		if t.stop {
			t.logger.Info("Stopping top row finder", zap.Int("jobType", t.jobType), zap.Int("tenant", t.tenant))
			break
		}
	}
}

func (q *queueImpl) Poll(ctx context.Context, req queue.PollRequest) (*queue.PollResponse, error) {
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
	partitionTime := endOfWeekPlusOneWeek(processAt)
	result = &queue.PollResponse{RecordPartitionTime: partitionTime}

	pollQuery := "SELECT id, pending_execution FROM jobs WHERE process_at=? AND tenant=? AND job_type=? AND state=? AND archive_after=? LIMIT 1 FOR UPDATE SKIP LOCKED"
	pollQuery = q.queryRewriter.RewriteQuery("jobs", pollQuery)
	err = tx.QueryRow(pollQuery, processAt, req.Tenant, req.JobType, queue.StatusScheduled, partitionTime).Scan(&result.Id, &remainingRetries)
	if err != nil {
		err = fmt.Errorf("failed to find top row for jobType=%d tenant=%d topTime=%s", req.JobType, req.Tenant, processAt.String())
		return nil, err
	}

	if remainingRetries <= 0 {
		if _, err = tx.Exec("UPDATE jobs SET state=?, sub_state=? WHERE id=?", queue.StatusFailed, queue.SubStatusNoRetryPendingError, result.Id); err != nil {
			return nil, fmt.Errorf("we found a record which is already dead: id=%s mark it failed and set it retries < 0", result.Id)
		}
	}

	// Update the row within the same transaction
	res, err := tx.Exec("UPDATE jobs SET state=?, version=version+1, pending_execution=pending_execution-1 WHERE id=? AND archive_after=?", queue.StatusProcessing, result.Id, partitionTime)
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
