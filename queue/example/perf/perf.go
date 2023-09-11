package perf

import (
	"context"
	"errors"
	"fmt"
	"github.com/devlibx/gox-base"
	queue "github.com/devlibx/gox-base/queue"
	mysqlQueue "github.com/devlibx/gox-base/queue/mysql"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var writeCounter, writeErrorCounter, readCounter, readErrorCounter, readWaitCounter, readNoResultCounter, fetchJobErrorCounter, markJobCompleteErrorCounter, markJobFailedErrorCounter metrics.Counter

var globalTenant = 2
var globalJobType = 1

var counter int32 = 0
var counterMutex = sync.Mutex{}
var startPerf = time.Now()
var fullStartPerf = time.Now()

func PerfMain() {
	argsWithoutProg := os.Args[1:]
	dontRunPoller := argsWithoutProg[0] == "w"

	writeCounter = metrics.NewCounter()
	readCounter = metrics.NewCounter()
	writeErrorCounter = metrics.NewCounter()
	readErrorCounter = metrics.NewCounter()
	readNoResultCounter = metrics.NewCounter()
	readWaitCounter = metrics.NewCounter()
	fetchJobErrorCounter = metrics.NewCounter()
	markJobCompleteErrorCounter = metrics.NewCounter()
	markJobFailedErrorCounter = metrics.NewCounter()
	metrics.Register("write", writeCounter)
	metrics.Register("write_error", writeErrorCounter)
	metrics.Register("read", readCounter)
	metrics.Register("read_error", readErrorCounter)
	metrics.Register("read_no_record", readNoResultCounter)
	metrics.Register("read_wait", readWaitCounter)
	metrics.Register("fetch_job_error", fetchJobErrorCounter)
	metrics.Register("mark_job_completed_error", markJobCompleteErrorCounter)
	metrics.Register("mark_job_failed_error", markJobFailedErrorCounter)

	go metrics.Log(metrics.DefaultRegistry, 1*time.Minute, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	// Setup DB
	storeBackend, err := mysqlQueue.NewMySqlBackedStore(queue.MySqlBackedStoreBackendConfig{
		Host:                 os.Getenv("DB_URL"),
		User:                 os.Getenv("DB_USER"),
		Password:             os.Getenv("DB_PASS"),
		Database:             os.Getenv("DB_NAME"),
		Port:                 3306,
		MaxOpenConnection:    100,
		MaxIdleConnection:    100,
		ConnMaxLifetimeInSec: 60,
		Properties:           gox.StringObjectMap{},
	}, true)
	if err != nil {
		panic(err)
	} else {
		db, _ := storeBackend.GetSqlDb()
		// db.Exec("TRUNCATE table jobs")
		// db.Exec("TRUNCATE table jobs_data")
		_ = db
	}

	// Setup ID generator - this will generate ULID based IDs
	idGenerator, err := queue.NewTimeBasedIdGenerator()
	if err != nil {
		panic(err)
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	crossFunction := gox.NewCrossFunction(zapConfig.Build())

	// Setup a queue to read/write
	appQueue, err := mysqlQueue.NewQueue(
		crossFunction,
		storeBackend,
		queue.MySqlBackedQueueConfig{
			Tenant:                     globalTenant,
			UsePreparedStatement:       true,
			UseMinQueryToPickLatestRow: true,
			DontRunPoller:              dontRunPoller,
		},
		idGenerator,
		queue.NewUdfAndTableNameQueryRewriter("jobs"),
	)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixMilli())

	// Run writer
	if argsWithoutProg[0] == "all" || argsWithoutProg[0] == "w" {
		go func() {
			for i := 0; i < 50; i++ {
				go func() {
					perfSchedule(appQueue)
				}()
			}
		}()
	}
	time.Sleep(2 * time.Second)

	// Run reader
	if argsWithoutProg[0] == "all" || argsWithoutProg[0] == "r" {
		go func() {
			for i := 0; i < 50; i++ {
				go func() {
					perfPoll(appQueue)
				}()
			}
		}()
	}

	time.Sleep(10000 * time.Hour)
}

var idSycn = sync.Mutex{}
var ids = map[string]string{}

func perfSchedule(appQueue queue.Queue) {
	var jobType = globalJobType
	var tenant = globalTenant

	for {
		now := time.Now()
		dayDelta := rand.Intn(30)
		hourDelta := rand.Intn(2)
		minDelta := rand.Intn(58)
		secondDelta := rand.Intn(58)

		// Add jobs to different tenant and to different job types
		jobType = rand.Intn(3)
		tenant = rand.Intn(3)

		// Generate records for all partitions
		ran := rand.Intn(5)
		if ran == 0 {
			now = time.Now()
		} else if ran == 1 {
			now = time.Now().Add(time.Duration(dayDelta*24) * time.Hour).Add(time.Duration(hourDelta) * time.Hour).Add(time.Duration(minDelta) * time.Minute).Add(time.Duration(secondDelta) * time.Second)
		} else if ran == 2 {
			now = time.Now().Add(time.Duration(-dayDelta*24) * time.Hour).Add(time.Duration(hourDelta) * time.Hour).Add(time.Duration(minDelta) * time.Minute).Add(time.Duration(secondDelta) * time.Second)
		} else {
			now = time.Now()
		}

		ctx, cancelContext := context.WithTimeout(context.Background(), 1*time.Second)
		if rs, err := appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            jobType,
			Tenant:             tenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", counter)},
		}); err != nil {
			fmt.Println("Failed to schedule job", err)
			writeErrorCounter.Inc(1)
			time.Sleep(1000 * time.Millisecond)
		} else {
			_ = rs
			writeCounter.Inc(1)
			time.Sleep(100 * time.Microsecond)

			idSycn.Lock()
			if r, ok := ids[rs.Id]; ok {
				fmt.Println("Clash", r)
				time.Sleep(10 * time.Second)
			} else {
				ids[rs.Id] = ""
			}
			idSycn.Unlock()

		}
		cancelContext()
	}
}

func perfPoll(appQueue queue.Queue) {
	for {
		ctx, cancelContext := context.WithTimeout(context.Background(), 1*time.Second)
		if rs, err := appQueue.Poll(ctx, queue.PollRequest{
			Tenant:  globalTenant,
			JobType: globalJobType,
		}); err != nil {
			var e *queue.PollResponseError
			if errors.As(err, &e) {
				readWaitCounter.Inc(1)
				time.Sleep(e.WaitForDurationBeforeTrying)
				continue
			} else {
				readErrorCounter.Inc(1)
				fmt.Println("Some error in polling", err)
			}
		} else {
			readCounter.Inc(1)
			c := atomic.AddInt32(&counter, 1)

			var fetchJobError, updateError error
			if true {
				fetchJobResponse, fetchJobErr := appQueue.FetchJobDetails(context.Background(), queue.JobDetailsRequest{Id: rs.Id})
				if fetchJobErr != nil {
					fetchJobError = fetchJobErr
					fetchJobErrorCounter.Inc(1)
				} else {
					if rand.Intn(5) == 0 && fetchJobErr == nil {
						// var result *queue.MarkJobFailedWithRetryResponse
						if _, updateErr := appQueue.MarkJobFailedAndScheduleRetry(
							context.Background(), queue.MarkJobFailedWithRetryRequest{Id: rs.Id, ScheduleRetryAt: fetchJobResponse.At.Add(time.Hour)},
						); updateErr != nil {
							markJobCompleteErrorCounter.Inc(1)
							updateError = updateErr
						}
					} else {
						if _, updateErr := appQueue.MarkJobCompleted(context.Background(), queue.MarkJobCompletedRequest{Id: rs.Id}); updateErr != nil {
							markJobFailedErrorCounter.Inc(1)
							updateError = updateErr
						}
					}
				}
			}

			// Log result
			if c%100 == 0 {
				counterMutex.Lock()
				fmt.Println(
					"Result =", rs,
					" count =", counter,
					" time =", time.Now().UnixMilli()-startPerf.UnixMilli(),
					"Total =", time.Now().UnixMilli()-fullStartPerf.UnixMilli(),
					"FetchJobError =", fetchJobError,
					"UpdateError =", updateError,
				)
				startPerf = time.Now()
				counterMutex.Unlock()
			}
		}
		cancelContext()
	}
}
