package queue

import (
	"context"
	"database/sql"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/devlibx/gox-base/serialization"
	"time"
)

func (q *queueImpl) jobInfoInit() (err error) {
	q.readJobDetailsOnce.Do(func() {
		jobQuery := "select job_type, state, sub_state, correlation_id, pending_execution FROM jobs WHERE id=? AND part=?"
		jobQuery = q.queryRewriter.RewriteQuery("jobs", jobQuery)
		jobDataQuery := "select properties, string_udf_1, string_udf_2, int_udf_1, int_udf_2 FROM jobs_data WHERE id=? AND part=?"
		jobDataQuery = q.queryRewriter.RewriteQuery("jobs_data", jobDataQuery)
		jobUpdateQuery := "UPDATE jobs set state=?, sub_state=? WHERE id=? AND part=?"
		jobUpdateQuery = q.queryRewriter.RewriteQuery("jobs", jobUpdateQuery)

		if q.readJobDetailsStatement, err = q.db.PrepareContext(context.Background(), jobQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for fetch job data")
		} else if q.readJobDataDetailsStatement, err = q.db.PrepareContext(context.Background(), jobDataQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for fetch job data")
		} else if q.updateJobStatusStatement, err = q.db.PrepareContext(context.Background(), jobUpdateQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for update job data")
		}
	})
	return
}

func (q *queueImpl) FetchJobDetails(ctx context.Context, req queue.JobDetailsRequest) (result *queue.JobDetailsResponse, err error) {
	return q.internalJobDetails(ctx, req)
}

func (q *queueImpl) internalJobDetails(ctx context.Context, req queue.JobDetailsRequest) (result *queue.JobDetailsResponse, err error) {
	result = &queue.JobDetailsResponse{}
	part := time.Time{}
	if err = q.jobInfoInit(); err != nil {
		return nil, errors.Wrap(err, "something is wrong we were not able to init read")
	} else if part, err = q.idToTime(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}
	part = partitionBasedOnProcessAtTime(part)

	var cid, strUdf1, strUdf2, properties sql.NullString
	var intUdf1, intUdf2 sql.NullInt64
	if err = q.readJobDetailsStatement.QueryRowContext(ctx, req.Id, part).Scan(&result.JobType, &result.State, &result.SubState, &cid, &result.RemainingExecution); err != nil {
		return nil, errors.Wrap(err, "failed to read job details: id=%s", req.Id)
	} else if err = q.readJobDataDetailsStatement.QueryRowContext(ctx, req.Id, part).Scan(&properties, &strUdf1, &strUdf2, &intUdf1, &intUdf2); err != nil {
		return nil, errors.Wrap(err, "failed to read job data details: id=%s", req.Id)
	}

	result.Id = req.Id
	result.At = part
	if cid.Valid {
		result.CorrelationId = cid.String
	}
	if strUdf1.Valid {
		result.StringUdf1 = strUdf1.String
	}
	if strUdf2.Valid {
		result.StringUdf2 = strUdf2.String
	}
	if cid.Valid {
		result.IntUdf1 = int(intUdf1.Int64)
	}
	if cid.Valid {
		result.IntUdf2 = int(intUdf1.Int64)
	}

	if properties.Valid {
		result.Properties = map[string]interface{}{}
		serialization.JsonBytesToObjectSuppressError([]byte(properties.String), &result.Properties)
	}

	return
}

func (q *queueImpl) MarkJobCompletedWithRetry(ctx context.Context, req queue.MarkJobFailedWithRetryRequest) (result *queue.MarkJobFailedWithRetryResponse, err error) {
	result = &queue.MarkJobFailedWithRetryResponse{Done: false}

	// Get the partition time
	part := time.Time{}
	if part, err = q.idToTime(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}
	part = partitionBasedOnProcessAtTime(part)

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
		}); err != nil {
			return nil, errors.Wrap(err, "failed to add new retry jobs (some retries are remaining for this job): id=%s", req.Id)
		}

		if _, err = q.updateJobStatusStatement.ExecContext(ctx, queue.StatusFailed, queue.SubStatusRetryPendingError, req.Id, part); err != nil {
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
	if part, err = q.idToTime(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}
	part = partitionBasedOnProcessAtTime(part)

	// Mark it done
	if _, err = q.updateJobStatusStatement.ExecContext(ctx, queue.StatusDone, queue.SubStatusDone, req.Id, part); err != nil {
		return nil, errors.Wrap(err, "failed to update the job: id=%s", req.Id)
	}
	return
}
