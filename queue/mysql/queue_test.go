package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/queue"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"testing"
	"time"
)

var dbHost = ""
var dbName = ""
var dbUser = ""
var dbPassword = ""
var testJobType = 79
var testTenant = 78

func setup() (storeBackend *mySqlStore, queueImpl queue.Queue, cf gox.CrossFunction, err error) {
	dbHost = os.Getenv("DB_URL")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")

	if storeBackend, err = NewMySqlBackedStore(queue.MySqlBackedStoreBackendConfig{
		Host:                 os.Getenv("DB_URL"),
		Port:                 3306,
		User:                 os.Getenv("DB_USER"),
		Password:             os.Getenv("DB_PASS"),
		Database:             os.Getenv("DB_NAME"),
		MaxOpenConnection:    100,
		MaxIdleConnection:    100,
		ConnMaxLifetimeInSec: 60,
		Properties:           gox.StringObjectMap{},
	}, true); err != nil {
		return
	}

	idGenerator, err := queue.NewTimeBasedIdGenerator()
	if err != nil {
		return
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	crossFunction := gox.NewCrossFunction(zapConfig.Build())
	if queueImpl, err = NewQueue(
		crossFunction,
		storeBackend,
		queue.MySqlBackedQueueConfig{
			Tenant:                     testTenant,
			UsePreparedStatement:       true,
			UseMinQueryToPickLatestRow: true,
		},
		idGenerator,
		queue.NewUdfAndTableNameQueryRewriter("jobs"),
	); err != nil {
		return
	}

	rand.Seed(time.Now().UnixMilli())
	return
}

func TestSchedule(t *testing.T) {
	if os.Getenv("DB_URL") == "" {
		t.Skip("to run tests you must set DB_URL which points to DB used in the test")
		return
	}

	now := time.Now()
	sc, appQueue, _, err := setup()
	assert.NoError(t, err)
	db := sc.db

	// Clear all test data if remaining
	markAllTestRowsToDone(t, context.Background(), db)

	t.Run("schedule a simple job", func(t *testing.T) {
		ctx, ch := context.WithTimeout(context.Background(), 10*time.Second)
		defer ch()

		rs, err := appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            testJobType,
			Tenant:             testTenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", 1023)},
		})
		assert.NoError(t, err)

		resultFromMySQL, readRowError := readRow(ctx, db, rs.Id)
		assert.NoError(t, readRowError)
		assert.Equal(t, queue.StatusScheduled, resultFromMySQL.State)
		assert.Equal(t, queue.SubStatusScheduledOk, resultFromMySQL.SubState)
		assert.Equal(t, 3, resultFromMySQL.RemainingExecution)
		assert.Equal(t, testTenant, resultFromMySQL.Tenant)
		assert.Equal(t, testJobType, resultFromMySQL.JobType)

		// Mark old test jobs completed - otherwise they just pile up
		_, _ = appQueue.MarkJobCompleted(ctx, queue.MarkJobCompletedRequest{Id: rs.Id})
	})

	t.Run("schedule a simple job and pull it back and also test if job mark complete works", func(t *testing.T) {
		ctx, ch := context.WithTimeout(context.Background(), 1*time.Second)
		defer ch()

		rs, err := appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            testJobType,
			Tenant:             testTenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", 1023)},
		})
		assert.NoError(t, err)

		resultFromMySQL, err := readRow(ctx, db, rs.Id)
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusScheduled, resultFromMySQL.State)
		assert.Equal(t, queue.SubStatusScheduledOk, resultFromMySQL.SubState)
		assert.Equal(t, 3, resultFromMySQL.RemainingExecution)
		assert.Equal(t, testTenant, resultFromMySQL.Tenant)
		assert.Equal(t, testJobType, resultFromMySQL.JobType)
		fmt.Println(rs)

		// Try to get the job - we pull 100 jobs in case we already have some jobs with test id
		foundJob := false
		var pollJobDetails *queue.JobDetailsResponse
		for i := 0; i < 100; i++ {
			pollResult, err := appQueue.Poll(ctx, queue.PollRequest{
				Tenant:  testTenant,
				JobType: testJobType,
			})
			assert.NoError(t, err, "since we just scheduled a job, we expect a job from the queue")

			// Job from some other test - ignore it
			if pollResult.Id != rs.Id {
				// Mark old test jobs completed - otherwise they just pile up
				_, _ = appQueue.MarkJobCompleted(ctx, queue.MarkJobCompletedRequest{Id: pollResult.Id})
				continue
			}
			foundJob = true

			pollJobDetails, err = appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: pollResult.Id})
			assert.NoError(t, err)

			assert.Equal(t, queue.StatusProcessing, pollJobDetails.State)
			assert.Equal(t, queue.SubStatusScheduledOk, pollJobDetails.SubState)
			assert.Equal(t, 2, pollJobDetails.RemainingExecution)
			assert.Equal(t, testTenant, pollJobDetails.Tenant)
			assert.Equal(t, testJobType, pollJobDetails.JobType)
			break
		}

		// Also test mark completed
		assert.True(t, foundJob)

		if pollJobDetails != nil {
			jobCompletedResponse, err := appQueue.MarkJobCompleted(ctx, queue.MarkJobCompletedRequest{Id: rs.Id})
			assert.NoError(t, err)
			assert.NotNil(t, jobCompletedResponse)
		}
	})
}

func TestPollWithNoRecord(t *testing.T) {
	if os.Getenv("DB_URL") == "" {
		t.Skip("to run tests you must set DB_URL which points to DB used in the test")
		return
	}

	sc, appQueue, _, err := setup()
	assert.NoError(t, err)

	// Clear all test data if remaining
	markAllTestRowsToDone(t, context.Background(), sc.db)

	t.Run("we expect no records", func(t *testing.T) {
		ctx, ch := context.WithTimeout(context.Background(), 10*time.Second)
		defer ch()

		now := time.Now()
		jobTime := now.Add(10 * time.Second)

		// Make all jobs completed for our test to tun
		db := sc.db
		_, err = db.ExecContext(context.Background(), "UPDATE jobs SET state=? WHERE tenant=? AND job_type=?", queue.StatusDone, testTenant, testJobType)
		assert.NoError(t, err)

		var scheduledResult *queue.ScheduleResponse
		scheduledResult, err = appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            testJobType,
			Tenant:             testTenant,
			At:                 jobTime,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", 1023)},
		})
		assert.NoError(t, err)
		_ = scheduledResult

		for i := 0; i < 10; i++ {
			_, err = appQueue.Poll(context.Background(), queue.PollRequest{
				Tenant:  testTenant,
				JobType: testJobType,
			})
			assert.Error(t, err, "since we just scheduled a job, we expect a job from the queue")

			var e *queue.PollResponseError
			tmp := errors.As(err, &e)
			assert.True(t, tmp)
			fmt.Println("Wait for ", e.WaitForDurationBeforeTrying.Milliseconds())
			assert.True(t, e.WaitForDurationBeforeTrying.Milliseconds() > 0)

			if tmp == true {
				break
			}
		}
	})

}

func TestRescheduleJobOnError(t *testing.T) {
	if os.Getenv("DB_URL") == "" {
		t.Skip("to run tests you must set DB_URL which points to DB used in the test")
		return
	}

	now := time.Now()
	sc, appQueue, _, err := setup()
	assert.NoError(t, err)
	db := sc.db

	t.Run("add retry on error", func(t *testing.T) {
		ctx, ch := context.WithTimeout(context.Background(), 10*time.Second)
		defer ch()

		var rs *queue.ScheduleResponse
		var pollResult *queue.PollResponse
		var resultFromMySQL, pollJobDetails *queue.JobDetailsResponse
		var jobFailedResponse *queue.MarkJobFailedWithRetryResponse
		var err error

		// Clear all test data if remaining
		markAllTestRowsToDone(t, ctx, db)

		// Part 1 - Schedule a job and verify its entry in DB
		rs, err = appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            testJobType,
			Tenant:             testTenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", 1023)},
		})
		assert.NoError(t, err)
		resultFromMySQL, err = readRow(ctx, db, rs.Id)
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusScheduled, resultFromMySQL.State)
		assert.Equal(t, queue.SubStatusScheduledOk, resultFromMySQL.SubState)
		assert.Equal(t, 3, resultFromMySQL.RemainingExecution)
		assert.Equal(t, testTenant, resultFromMySQL.Tenant)
		assert.Equal(t, testJobType, resultFromMySQL.JobType)
		fmt.Println(rs)

		// Part 2 - get the job and also fetch its complete details
		pollResult, err = appQueue.Poll(ctx, queue.PollRequest{
			Tenant:  testTenant,
			JobType: testJobType,
		})
		assert.NoError(t, err, "since we just scheduled a job, we expect a job from the queue")

		// Part 2.1 - fetch its complete details
		pollJobDetails, err = appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: pollResult.Id})
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusProcessing, pollJobDetails.State)
		assert.Equal(t, queue.SubStatusScheduledOk, pollJobDetails.SubState)
		assert.Equal(t, 2, pollJobDetails.RemainingExecution)
		assert.Equal(t, testTenant, pollJobDetails.Tenant)
		assert.Equal(t, testJobType, pollJobDetails.JobType)

		// Part 3 - Mark it failed
		jobFailedResponse, err = appQueue.MarkJobFailedAndScheduleRetry(ctx, queue.MarkJobFailedWithRetryRequest{
			Id:              pollResult.Id,
			ScheduleRetryAt: now.Add(10 * time.Millisecond),
		})
		time.Sleep(10 * time.Millisecond)
		assert.NoError(t, err)
		newJobId := jobFailedResponse.RetryJobId

		// Part 3.1 - make sure we marked the original job as failed
		jd, err := readRow(ctx, db, pollResult.Id)
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusFailed, jd.State)
		assert.Equal(t, queue.SubStatusRetryPendingError, jd.SubState)
		assert.Equal(t, 2, jd.RemainingExecution)
		assert.Equal(t, testTenant, jd.Tenant)
		assert.Equal(t, testJobType, jd.JobType)

		// Part 4- Get the job again and this time it should be the retry job
		// Also the poll result should have the id of the new job id
		pollResult, err = appQueue.Poll(ctx, queue.PollRequest{
			Tenant:  testTenant,
			JobType: testJobType,
		})
		assert.NoError(t, err, "since we just scheduled a job, we expect a job from the queue")
		assert.Equal(t, newJobId, pollResult.Id)

		// Part 2.1 - fetch its complete details
		pollJobDetails, err = appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: pollResult.Id})
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusProcessing, pollJobDetails.State)
		assert.Equal(t, queue.SubStatusScheduledOk, pollJobDetails.SubState)
		assert.Equal(t, 1, pollJobDetails.RemainingExecution)
		assert.Equal(t, testTenant, pollJobDetails.Tenant)
		assert.Equal(t, testJobType, pollJobDetails.JobType)

		// Again fail it
		jobFailedResponse, err = appQueue.MarkJobFailedAndScheduleRetry(ctx, queue.MarkJobFailedWithRetryRequest{
			Id:              pollResult.Id,
			ScheduleRetryAt: now.Add(10 * time.Millisecond),
		})
		time.Sleep(10 * time.Millisecond)
		assert.NoError(t, err)
		newJobId = jobFailedResponse.RetryJobId
		pollResult, err = appQueue.Poll(ctx, queue.PollRequest{
			Tenant:  testTenant,
			JobType: testJobType,
		})
		assert.NoError(t, err, "since we just scheduled a job, we expect a job from the queue")
		assert.Equal(t, newJobId, pollResult.Id)
		pollJobDetails, err = appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: pollResult.Id})
		assert.NoError(t, err)
		assert.Equal(t, 0, pollJobDetails.RemainingExecution)
		jobFailedResponse, err = appQueue.MarkJobFailedAndScheduleRetry(ctx, queue.MarkJobFailedWithRetryRequest{
			Id:              pollResult.Id,
			ScheduleRetryAt: now.Add(10 * time.Millisecond),
		})
		time.Sleep(10 * time.Millisecond)
		assert.NoError(t, err)
		pollJobDetails, err = appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: pollResult.Id})
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusFailed, pollJobDetails.State)
		assert.Equal(t, queue.SubStatusNoRetryPendingError, pollJobDetails.SubState)

		var count int
		err = db.QueryRowContext(ctx, "SELECT count(*) FROM jobs_data WHERE retry_group=?", pollJobDetails.RetryGroup).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
		fmt.Println(pollJobDetails.RetryGroup)

		fmt.Println(newJobId)
	})

}

func TestUpdateJobData(t *testing.T) {
	if os.Getenv("DB_URL") == "" {
		t.Skip("to run tests you must set DB_URL which points to DB used in the test")
		return
	}

	sc, appQueue, _, err := setup()
	assert.NoError(t, err)
	// db := sc.db
	now := time.Now()

	// Clear all test data if remaining
	markAllTestRowsToDone(t, context.Background(), sc.db)

	t.Run("update job data", func(t *testing.T) {
		ctx, ch := context.WithTimeout(context.Background(), 10*time.Second)
		defer ch()

		// Part 1 - Schedule a job and verify its entry in DB
		rs, err := appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            testJobType,
			Tenant:             testTenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", 5510)},
			StringUdf1:         "str_udf_1",
			StringUdf2:         "str_udf_2",
			IntUdf1:            10,
			IntUdf2:            11,
		})
		assert.NoError(t, err)
		resultFromMySQL, err := appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: rs.Id})
		assert.NoError(t, err)
		assert.Equal(t, queue.StatusScheduled, resultFromMySQL.State)
		assert.Equal(t, queue.SubStatusScheduledOk, resultFromMySQL.SubState)
		assert.Equal(t, 3, resultFromMySQL.RemainingExecution)
		assert.Equal(t, testTenant, resultFromMySQL.Tenant)
		assert.Equal(t, testJobType, resultFromMySQL.JobType)
		assert.Equal(t, "str_udf_1", resultFromMySQL.StringUdf1)
		assert.Equal(t, "str_udf_2", resultFromMySQL.StringUdf2)
		assert.Equal(t, 10, resultFromMySQL.IntUdf1)
		assert.Equal(t, 11, resultFromMySQL.IntUdf2)
		assert.Equal(t, "5510", resultFromMySQL.Properties["info"])
		fmt.Println(rs)

		updateResponse, err := appQueue.UpdateJobData(ctx, queue.UpdateJobDataRequest{
			Id:         rs.Id,
			StringUdf1: "str_udf_1_updated",
			StringUdf2: "str_udf_2_updated",
			IntUdf1:    110,
			IntUdf2:    111,
			Properties: map[string]interface{}{"info": fmt.Sprintf("%d", 15510)},
		})
		assert.NoError(t, err)
		assert.NotNil(t, updateResponse)

		jd, err := appQueue.FetchJobDetails(ctx, queue.JobDetailsRequest{Id: rs.Id})
		assert.NoError(t, err)
		assert.NotNil(t, jd)
		assert.Equal(t, queue.SubStatusScheduledOk, jd.SubState)
		assert.Equal(t, 3, jd.RemainingExecution)
		assert.Equal(t, testTenant, jd.Tenant)
		assert.Equal(t, testJobType, jd.JobType)
		assert.Equal(t, "str_udf_1_updated", jd.StringUdf1)
		assert.Equal(t, "str_udf_2_updated", jd.StringUdf2)
		assert.Equal(t, 110, jd.IntUdf1)
		assert.Equal(t, 111, jd.IntUdf2)
		assert.Equal(t, "15510", jd.Properties["info"])
	})
}

func readRow(ctx context.Context, db *sql.DB, id string) (result *queue.JobDetailsResponse, err error) {
	result = &queue.JobDetailsResponse{}

	err = db.QueryRowContext(ctx, "SELECT job_type, pending_execution, tenant, state, sub_state FROM jobs WHERE id=?", id).
		Scan(
			&result.JobType,
			&result.RemainingExecution,
			&result.Tenant,
			&result.State,
			&result.SubState,
		)
	return
}

func markAllTestRowsToDone(t *testing.T, ctx context.Context, db *sql.DB) {
	_, err := db.ExecContext(ctx, "UPDATE jobs SET state=? WHERE tenant=? AND job_type=?", queue.StatusDone, testTenant, testJobType)
	assert.NoError(t, err)
}
