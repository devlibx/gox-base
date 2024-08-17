package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/devlibx/gox-base/v2/queue"
	"github.com/devlibx/gox-base/v2/serialization"
	"go.uber.org/zap"
	"time"
)

func (q *queueImpl) MarkJobFailedAndScheduleRetry(ctx context.Context, req queue.MarkJobFailedWithRetryRequest) (result *queue.MarkJobFailedWithRetryResponse, err error) {
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
		result.Done = false
	} else if req.ScheduleRetryAt.IsZero() {
		if _, err = q.updateJobStatusStatement.ExecContext(ctx, queue.StatusFailed, queue.SubStatusRetryIgnoredByUserError, req.Id, part); err != nil {
			return nil, errors.Wrap(err, "failed to update the job to mark failed: id=%s", req.Id)
		}
		result.Done = false
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
			At:                   req.ScheduleRetryAt,
			JobType:              jobFetchResponse.JobType,
			Tenant:               jobFetchResponse.Tenant,
			CorrelationId:        jobFetchResponse.CorrelationId,
			RemainingExecution:   jobFetchResponse.RemainingExecution,
			StringUdf1:           jobFetchResponse.StringUdf1,
			StringUdf2:           jobFetchResponse.StringUdf2,
			IntUdf1:              jobFetchResponse.IntUdf1,
			IntUdf2:              jobFetchResponse.IntUdf1,
			Properties:           jobFetchResponse.Properties,
			InternalRetryGroupId: jobFetchResponse.RetryGroup,
			InternalTx:           tx,
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

func (q *queueImpl) UpdateJobData(ctx context.Context, req queue.UpdateJobDataRequest) (result *queue.UpdateJobDataResponse, err error) {
	return q.internalUpdateJobData(ctx, req)
}

func (q *queueImpl) internalUpdateJobData(ctx context.Context, req queue.UpdateJobDataRequest) (result *queue.UpdateJobDataResponse, err error) {
	result = &queue.UpdateJobDataResponse{}

	// Get the partition time
	part := time.Time{}
	if part, err = queue.GeneratePartitionTimeByRecordId(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}

	properties := ""
	if req.Properties != nil {
		// properties = req.Properties
		if properties, err = serialization.Stringify(req.Properties); err != nil {
			return nil, fmt.Errorf("failed to persist (metadata is bad): %w", err)
		}
	}

	// Update job data
	var r sql.Result
	var noOfUpdatedRecords int64
	if r, err = q.updateJobDataStatement.ExecContext(ctx, req.StringUdf1, req.StringUdf2, req.IntUdf1, req.IntUdf2, properties, req.Id, part); err != nil {
		return nil, errors.Wrap(err, "failed to update the job data: id=%s", req.Id)
	} else if noOfUpdatedRecords, err = r.RowsAffected(); err == nil && noOfUpdatedRecords == 0 {
		err = errors.Wrap(err, "failed to update the job data : id=%s", req.Id)
	}

	return
}
