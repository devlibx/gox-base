package queue

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/devlibx/gox-base"
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
		Host:          os.Getenv("DB_URL"),
		Port:          3306,
		User:          os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASS"),
		Database:      os.Getenv("DB_NAME"),
		MaxConnection: 100,
		MinConnection: 100,
		Properties:    gox.StringObjectMap{},
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
