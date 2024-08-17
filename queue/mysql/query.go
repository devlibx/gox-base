package queue

import (
	"context"
	"database/sql"
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/devlibx/gox-base/v2/queue"
	"github.com/devlibx/gox-base/v2/serialization"
	"time"
)

func (q *queueImpl) jobInfoInit() (err error) {
	q.readJobDetailsOnce.Do(func() {
		jobQuery := "select job_type, state, sub_state, correlation_id, pending_execution, tenant FROM jobs WHERE id=? AND part=?"
		jobQuery = q.queryRewriter.RewriteQuery("jobs", jobQuery)
		jobDataQuery := "select properties, string_udf_1, string_udf_2, int_udf_1, int_udf_2, retry_group FROM jobs_data WHERE id=? AND part=?"
		jobDataQuery = q.queryRewriter.RewriteQuery("jobs_data", jobDataQuery)
		jobUpdateQuery := "UPDATE jobs set state=?, sub_state=? WHERE id=? AND part=?"
		jobUpdateQuery = q.queryRewriter.RewriteQuery("jobs", jobUpdateQuery)
		jobDataUpdateQuery := "UPDATE jobs_data SET string_udf_1=?, string_udf_2=?, int_udf_1=?, int_udf_2=?, properties=? WHERE id=? AND part=?"
		jobDataUpdateQuery = q.queryRewriter.RewriteQuery("jobs_data", jobDataUpdateQuery)

		if q.readJobDetailsStatement, err = q.db.PrepareContext(context.Background(), jobQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for fetch job data")
		} else if q.readJobDataDetailsStatement, err = q.db.PrepareContext(context.Background(), jobDataQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for fetch job data")
		} else if q.updateJobStatusStatement, err = q.db.PrepareContext(context.Background(), jobUpdateQuery); err != nil {
			err = errors.Wrap(err, "failed to build query for update job status")
		} else if q.updateJobDataStatement, err = q.db.PrepareContext(context.Background(), jobDataUpdateQuery); err != nil {
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
	} else if part, err = queue.GeneratePartitionTimeByRecordId(req.Id); err != nil {
		return nil, errors.Wrap(err, "not able to get time out of id: id=%s", req.Id)
	}

	var cid, strUdf1, strUdf2, properties, retryGroup sql.NullString
	var intUdf1, intUdf2 sql.NullInt64
	var tenant sql.NullInt32
	if err = q.readJobDetailsStatement.QueryRowContext(ctx, req.Id, part).Scan(&result.JobType, &result.State, &result.SubState, &cid, &result.RemainingExecution, &tenant); err != nil {
		return nil, errors.Wrap(err, "failed to read job details: id=%s", req.Id)
	} else if err = q.readJobDataDetailsStatement.QueryRowContext(ctx, req.Id, part).Scan(&properties, &strUdf1, &strUdf2, &intUdf1, &intUdf2, &retryGroup); err != nil {
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
		result.IntUdf2 = int(intUdf2.Int64)
	}
	if tenant.Valid {
		result.Tenant = int(tenant.Int32)
	}
	if retryGroup.Valid {
		result.RetryGroup = retryGroup.String
	}

	if properties.Valid {
		result.Properties = map[string]interface{}{}
		serialization.JsonBytesToObjectSuppressError([]byte(properties.String), &result.Properties)
	}

	return
}
