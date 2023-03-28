package goxServer

import (
	"context"
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"os"
	"testing"
	"time"
)

type testInternal struct {
}

func (t *testInternal) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.Write([]byte(`{"status": "ok"}`))
	w.WriteHeader(http.StatusOK)
}

func TestGoxServer(t *testing.T) {
	if os.Getenv("TEST_GOX_SERVER") == "" {
		// t.Skip("Skip gox server test")
	}

	var loggerConfig zap.Config
	loggerConfig = zap.NewDevelopmentConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, _ := loggerConfig.Build()

	cf := gox.NewNoOpCrossFunction(logger)
	server, err := NewServer(cf)
	assert.NoError(t, err)

	err = server.Start(&testInternal{}, &config.App{
		AppName:                     "test",
		HttpPort:                    15678,
		Environment:                 "tst",
		RequestReadTimeoutMs:        10000,
		RequestWriteTimeoutMs:       10000,
		OutstandingRequestTimeoutMs: 10000,
		IdleTimeoutMs:               1000,
	})

	go func() {
		err = server.Start(&testInternal{}, &config.App{
			AppName:                     "test",
			HttpPort:                    15678,
			Environment:                 "tst",
			RequestReadTimeoutMs:        10000,
			RequestWriteTimeoutMs:       10000,
			OutstandingRequestTimeoutMs: 10000,
			IdleTimeoutMs:               1000,
		})
	}()
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)

	writerClient := resty.New()
	writerClient.SetBaseURL("http://localhost:5678")

	// Setup 1 - Store new tenant definition
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	resp, err := writerClient.R().
		SetContext(ctx).
		Post("/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	fmt.Println(string(resp.Body()))
	if err != nil {
		panic(err)
	}
}
