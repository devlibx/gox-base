package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	queue "github.com/devlibx/gox-base/queue"
	mysqlQueue "github.com/devlibx/gox-base/queue/mysql"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"time"
)

func main() {
	storeBackend, err := mysqlQueue.NewMySqlBackedStore(queue.MySqlBackedStoreBackendConfig{
		Host:          os.Getenv("DB_URL"),
		Port:          3306,
		User:          os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASS"),
		Database:      os.Getenv("DB_NAME"),
		MaxConnection: 50,
		MinConnection: 50,
		Properties:    gox.StringObjectMap{},
	}, true)
	if err != nil {
		panic(err)
	} else {
		db, _ := storeBackend.GetSqlDb()
		db.Exec("TRUNCATE table jobs")
		db.Exec("TRUNCATE table jobs_data")
	}

	idGenerator, err := queue.NewTimeBasedIdGenerator()
	// idGenerator, err := queue.NewRandomUuidIdGenerator()
	if err != nil {
		panic(err)
	}

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	crossFunction := gox.NewCrossFunction(zapConfig.Build())

	appQueue, err := mysqlQueue.NewQueue(
		crossFunction,
		storeBackend,
		queue.MySqlBackedQueueConfig{
			Tenant:     11,
			MaxJobType: 1,
		},
		idGenerator,
		queue.NewUdfAndTableNameQueryRewriter("jobs"),
	)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixMilli())
	go func() {
		for i := 0; i < 10; i++ {
			go func() {
				push(appQueue)
			}()
		}

	}()

	go func() {
		for i := 0; i < 10; i++ {
			go func() {
				poll(appQueue)
			}()
		}

	}()

	time.Sleep(time.Hour)
}

func push(appQueue queue.Queue) {
	if true {
		// return
	}
	count := 0

	for {
		start := time.Now()
		h := rand.Intn(200)
		m := rand.Intn(50)
		_, _, _ = start, h, m
		now := time.Now().Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Second)
		now = time.Now()
		rs, err := appQueue.Schedule(context.Background(), queue.ScheduleRequest{
			JobType:    1,
			Tenant:     11,
			At:         now,
			Properties: map[string]interface{}{"info": fmt.Sprintf("%d", count)},
		})
		if err != nil {
			fmt.Println("Error", err)
			time.Sleep(1 * time.Second)
		} else {
			// fmt.Printf("Result = Ok: Id=%-30s  Hour=%-3d Min=%-3d  TimeTakne=%-4d \n", rs.Id, h, m, time.Now().UnixMilli()-start.UnixMilli())
			time.Sleep(10 * time.Microsecond)
		}
		_ = rs
		count++

		if count > 3000 {
			break
		}
	}
}

func poll(appQueue queue.Queue) {
	for {
		rs, err := appQueue.Poll(context.Background(), queue.PollRequest{
			Tenant:  11,
			JobType: 1,
		})
		if err != nil {
			fmt.Println("Error", err)
			time.Sleep(1 * time.Second)
		} else {
			// fmt.Printf("Result = Ok: Id=%-30s  Hour=%-3d Min=%-3d  TimeTakne=%-4d \n", rs.Id, h, m, time.Now().UnixMilli()-start.UnixMilli())
			fmt.Println("Result Ok", rs)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
