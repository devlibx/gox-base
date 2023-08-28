package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	queue "github.com/devlibx/gox-base/queue"
	mysqlQueue "github.com/devlibx/gox-base/queue/mysql"
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
	}

	idGenerator, err := queue.NewTimeBasedIdGenerator()
	// idGenerator, err := queue.NewRandomUuidIdGenerator()
	if err != nil {
		panic(err)
	}

	appQueue, err := mysqlQueue.NewQueue(
		storeBackend,
		queue.MySqlBackedQueueConfig{},
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

	time.Sleep(time.Hour)
}

func push(appQueue queue.Queue) {
	count := 0

	for {
		start := time.Now()
		h := rand.Intn(200)
		m := rand.Intn(50)
		now := time.Now().Add(time.Duration(h) * time.Hour).Add(time.Duration(m) * time.Second)
		rs, err := appQueue.Schedule(context.Background(), queue.ScheduleRequest{
			JobType:    "cron",
			At:         now,
			Properties: map[string]interface{}{"info": fmt.Sprintf("%d", count)},
		})
		if err != nil {
			fmt.Println("Error", err)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Ok", rs.Id, " h=", h, " m=", m, "       TimeTaken=", (time.Now().UnixMilli() - start.UnixMilli()))
			time.Sleep(100 * time.Microsecond)
		}
		_ = rs
		count++
	}
}
