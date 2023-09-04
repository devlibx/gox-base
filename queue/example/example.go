package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/devlibx/gox-base"
	queue "github.com/devlibx/gox-base/queue"
	"github.com/devlibx/gox-base/queue/example/perf"
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

// r := NewRegistry()

func main() {
	if true {
		perf.PerfMain()
		return
	}

	//	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

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

	storeBackend, err := mysqlQueue.NewMySqlBackedStore(queue.MySqlBackedStoreBackendConfig{
		Host:                 os.Getenv("DB_URL"),
		Port:                 3306,
		User:                 os.Getenv("DB_USER"),
		Password:             os.Getenv("DB_PASS"),
		Database:             os.Getenv("DB_NAME"),
		MaxOpenConnection:    100,
		MaxIdleConnection:    100,
		ConnMaxLifetimeInSec: 60,
		Properties:           gox.StringObjectMap{},
	}, true)
	if err != nil {
		panic(err)
	} else {
		db, _ := storeBackend.GetSqlDb()
		db.Exec("TRUNCATE table jobs")
		db.Exec("TRUNCATE table jobs_data")
		_ = db
	}

	idGenerator, err := queue.NewTimeBasedIdGenerator()
	// idGenerator, err := queue.NewRandomUuidIdGenerator()
	if err != nil {
		panic(err)
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	crossFunction := gox.NewCrossFunction(zapConfig.Build())

	dontRunPoller := argsWithoutProg[0] == "w"

	appQueue, err := mysqlQueue.NewQueue(
		crossFunction,
		storeBackend,
		queue.MySqlBackedQueueConfig{
			Tenant:                     globalTenant,
			MaxJobType:                 1,
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

	if argsWithoutProg[0] == "all" || argsWithoutProg[0] == "w" {
		go func() {
			for i := 0; i < 50; i++ {
				go func() {
					push(appQueue)
				}()
			}

		}()
	}

	time.Sleep(2 * time.Second)
	if argsWithoutProg[0] == "all" || argsWithoutProg[0] == "r" {
		go func() {
			for i := 0; i < 50; i++ {
				go func() {
					poll(appQueue)
				}()
			}

		}()
	}

	time.Sleep(time.Hour)
}

var rowInsertedCount int64 = 0

func push(appQueue queue.Queue) {
	if true {
		// return
	}

	for {
		start := time.Now()
		d := rand.Intn(30)
		h := rand.Intn(2)
		m := rand.Intn(58)
		s := rand.Intn(58)
		_, _, _, _, _ = start, h, m, s, d
		var now time.Time
		if rand.Int()%5 == 0 {
			now = time.Now()
		} else {
			now = time.Now().Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Minute).Add(time.Duration(s) * time.Second)
		}

		var jobType = globalJobType
		var tenant = globalTenant
		/*
			now = time.Now()

			ran := rand.Intn(5)
			if ran == 0 {
				now = time.Now()
			} else if ran == 1 {
				now = time.Now().Add(time.Duration(d*24) * time.Hour).Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Minute).Add(time.Duration(s) * time.Second)
			} else if ran == 2 {
				now = time.Now().Add(time.Duration(-d*24) * time.Hour).Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Minute).Add(time.Duration(s) * time.Second)
			} else {
				now = time.Now()
			}

			jobType := rand.Intn(3)
			tenant := rand.Intn(3)


					if rand.Intn(2)%2 == 0 {
						now = time.Now()
						jobType = globalJobType
						tenant = globalTenant
					}
				if true {
					now = time.Now()
					jobType = globalJobType
					tenant = globalTenant
				}

		*/

		ctx, ch := context.WithTimeout(context.Background(), 1*time.Second)
		rs, err := appQueue.Schedule(ctx, queue.ScheduleRequest{
			JobType:            jobType,
			Tenant:             tenant,
			At:                 now,
			RemainingExecution: 3,
			Properties:         map[string]interface{}{"info": fmt.Sprintf("%d", count)},
		})
		ch()
		if err != nil {
			fmt.Println("Error", err)
			writeErrorCounter.Inc(1)
			time.Sleep(1000 * time.Millisecond)
		} else {
			// fmt.Printf("Result = Ok: Id=%-30s  Hour=%-3d Min=%-3d  TimeTakne=%-4d \n", rs.Id, h, m, time.Now().UnixMilli()-start.UnixMilli())
			writeCounter.Inc(1)
			time.Sleep(100 * time.Microsecond)
		}
		_ = rs

		if atomic.AddInt64(&rowInsertedCount, 1) > 30000000000 {
			break
		}

		time.Sleep(1 * time.Millisecond)
	}
}

var count int32 = 0
var m = sync.Mutex{}
var start = time.Now()
var fullStart = time.Now()
var ma = map[string]string{}
var cM = sync.Mutex{}

func poll(appQueue queue.Queue) {
	for {
		ctx, ch := context.WithTimeout(context.Background(), 1*time.Second)
		rs, err := appQueue.Poll(ctx, queue.PollRequest{
			Tenant:  globalTenant,
			JobType: globalJobType,
		})
		ch()
		_ = rs
		if err != nil {
			var e *queue.PollResponseError
			if errors.As(err, &e) {
				fmt.Println("Wait for sometime", e)
				writeErrorCounter.Inc(1)
				time.Sleep(e.WaitForDurationBeforeTrying)
				continue
			} else {
				fmt.Println("Some error in polling", err)
			}
		} else {
			// fmt.Printf("Result = Ok: Id=%-30s  Hour=%-3d Min=%-3d  TimeTakne=%-4d \n", rs.Id, h, m, time.Now().UnixMilli()-start.UnixMilli())
			readCounter.Inc(1)
			c := atomic.AddInt32(&count, 1)
			/*
				if c%100 == 0 {
					m.Lock()
					// fmt.Println("Result Ok", rs, " count=", c, " time=", time.Now().UnixMilli()-start.UnixMilli(), "      Total=", time.Now().UnixMilli()-fullStart.UnixMilli())
					start = time.Now()
					m.Unlock()
				}
			*/

			// fmt.Println("Result Ok", rs, " count=", count)

			cM.Lock()
			if k, ok := ma[rs.Id]; ok {
				fmt.Println("--> dup ", k)
			}
			ma[rs.Id] = "a"
			cM.Unlock()

			time.Sleep(100 * time.Microsecond)

			var updateErr, getErr error
			var jd *queue.JobDetailsResponse
			_ = jd

			if true {
				jd, getErr = appQueue.FetchJobDetails(context.Background(), queue.JobDetailsRequest{Id: rs.Id})
				if getErr != nil {
					fmt.Println("failed to fetch job details")
				} else {
					// fmt.Println("Job details", jobD)
				}

				if rand.Intn(5) == 0 && getErr == nil {
					var result *queue.MarkJobFailedWithRetryResponse
					if result, updateErr = appQueue.MarkJobFailedAndScheduleRetry(context.Background(), queue.MarkJobFailedWithRetryRequest{Id: rs.Id, ScheduleRetryAt: jd.At.Add(time.Hour)}); updateErr != nil {
						fmt.Println("failed to marked job failed", err)
					} else {
						// fmt.Println("OK failed to marked job failed", result.RetryJobId)
					}
					_ = result
				} else {
					if _, updateErr = appQueue.MarkJobCompleted(context.Background(), queue.MarkJobCompletedRequest{Id: rs.Id}); updateErr != nil {
						// fmt.Println("failed to mark job done", rs.Id)
					}
				}

			}

			if c%100 == 0 {
				m.Lock()
				// fmt.Println("Result Ok", rs, " count=", c, " time=", time.Now().UnixMilli()-start.UnixMilli(), "      Total=", time.Now().UnixMilli()-fullStart.UnixMilli())
				fmt.Println("Result Ok", rs, " count=", count, "GetErr=", getErr, " UpdateError=", updateErr, " time=", time.Now().UnixMilli()-start.UnixMilli(), "      Total=", time.Now().UnixMilli()-fullStart.UnixMilli())
				start = time.Now()
				m.Unlock()
			}

		}
	}
}
