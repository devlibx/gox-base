package main

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	queue "github.com/devlibx/gox-base/queue"
	mysqlQueue "github.com/devlibx/gox-base/queue/mysql"
	"os"
	"time"
)

func main() {
	storeBackend, err := mysqlQueue.NewMySqlBackedStore(queue.MySqlBackedStoreBackendConfig{
		Host:          "localhost",
		Port:          3306,
		User:          "root",
		Password:      os.Getenv("DB_PASS"),
		Database:      "users",
		MaxConnection: 50,
		MinConnection: 50,
		Properties:    gox.StringObjectMap{},
	}, true)
	if err != nil {
		panic(err)
	}

	idGenerator, err := queue.NewTimeBasedIdGenerator()
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

	go func() {
		push(appQueue)
	}()

	time.Sleep(time.Hour)
}

func push(appQueue queue.Queue) {
	count := 0
	for {
		now := time.Now()
		rs, err := appQueue.Schedule(context.Background(), queue.ScheduleRequest{
			JobType:    "cron",
			At:         now,
			Properties: map[string]interface{}{"info": fmt.Sprintf("%d", count)},
		})
		if err != nil {
			fmt.Println("Error", err)
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Ok", rs.Id)
			time.Sleep(1 * time.Millisecond)
		}
		_ = rs
		count++
	}
}
