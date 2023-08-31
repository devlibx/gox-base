package queue

import (
	"context"
	"database/sql"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"go.uber.org/zap"
	"time"
)

func (q *queueImpl) MarkJobCompletedWithRetry(ctx context.Context, req queue.MarkJobFailedWithRetryRequest) (result *queue.MarkJobFailedWithRetryResponse, err error) {
	result = &queue.MarkJobFailedWithRetryResponse{Done: false}

	// Get the partition time
	part := time.Time{}
	if part, err = queue.GeneratePartitionTimeByRecordId(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}

	var jobFetchResponse *queue.JobDetailsResponse
	if jobFetchResponse, err = q.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: req.Id}); err != nil {
		return nil, errors.Wrap(err, "failed to get job with id=%s - needed to setup retry", req.Id)
	}

	if jobFetchResponse.RemainingExecution <= 0 {
		if _, err = q.updateJobStatusStatement.ExecContext(ctx, queue.StatusFailed, queue.SubStatusNoRetryPendingError, req.Id, part); err != nil {
			return nil, errors.Wrap(err, "failed to update the job to mark failed: id=%s", req.Id)
		}
		result.Done = true
	} else {

		// Begin a transaction
		var tx *sql.Tx
		tx, err = q.db.Begin()
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

		var scheduleResponse *queue.ScheduleResponse
		if scheduleResponse, err = q.Schedule(ctx, queue.ScheduleRequest{
			At:                 req.ScheduleRetryAt,
			JobType:            jobFetchResponse.JobType,
			Tenant:             jobFetchResponse.Tenant,
			CorrelationId:      jobFetchResponse.CorrelationId,
			RemainingExecution: jobFetchResponse.RemainingExecution - 1,
			StringUdf1:         jobFetchResponse.StringUdf1,
			StringUdf2:         jobFetchResponse.StringUdf2,
			IntUdf1:            jobFetchResponse.IntUdf1,
			IntUdf2:            jobFetchResponse.IntUdf1,
			Properties:         jobFetchResponse.Properties,
			InternalTx:         tx,
		}); err != nil {
			return nil, errors.Wrap(err, "failed to add new retry jobs (some retries are remaining for this job): id=%s", req.Id)
		}

		if _, err = tx.StmtContext(ctx, q.updateJobStatusStatement).ExecContext(ctx, queue.StatusFailed, queue.SubStatusRetryPendingError, req.Id, part); err != nil {
			return nil, errors.Wrap(err, "failed to update the job to mark failed: id=%s", req.Id)
		}

		result.RetryJobId = scheduleResponse.Id
		result.Done = true
	}

	return
}

func (q *queueImpl) MarkJobCompleted(ctx context.Context, req queue.MarkJobCompletedRequest) (result *queue.MarkJobCompletedResponse, err error) {
	result = &queue.MarkJobCompletedResponse{}

	// Get the partition time
	part := time.Time{}
	if part, err = queue.GeneratePartitionTimeByRecordId(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}

	// Mark it done
	if _, err = q.updateJobStatusStatement.ExecContext(ctx, queue.StatusDone, queue.SubStatusDone, req.Id, part); err != nil {
		return nil, errors.Wrap(err, "failed to update the job: id=%s", req.Id)
	}
	return
}
