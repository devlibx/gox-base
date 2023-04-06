package goxSql

import (
	"context"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/util"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var regexToCleanQueryToDump = regexp.MustCompile(`^[^\n]+\n`)
var startMetricDumpSyncOnce = sync.Once{}

type LogInfo struct {
	ctx                         context.Context
	name                        string
	startTime                   int64
	timeTaken                   int64
	query                       string
	cleanQuery                  string
	logger                      *zap.Logger
	hist                        metrics.Histogram
	callbacks                   *Callbacks
	enableSqlQueryLogging       bool
	enableSqlQueryMetricLogging bool
}

// Done is the method which is normally set in "defer". This will dump the required logs, post callbacks etc
func (l LogInfo) Done(err error, args ...interface{}) {
	endTime := time.Now().UnixMilli()
	l.timeTaken = endTime - l.startTime

	if l.enableSqlQueryMetricLogging {
		q := l.cleanQuery
		l.hist = metrics.GetOrRegisterHistogram(q, metrics.DefaultRegistry, metrics.NewExpDecaySample(1028, 0.015))
		l.hist.Update(l.timeTaken)
	}

	if l.enableSqlQueryLogging {
		l.logger.Info(l.name, zap.Int64("time", l.timeTaken), zap.String("query", l.cleanQuery), zap.Any("args", args))
	}

	if l.callbacks != nil {
		l.callbacks.PostCallbackFunc(PostCallbackData{
			Ctx:       l.ctx,
			Name:      l.name,
			StartTime: l.startTime,
			EndTime:   endTime,
			TimeTaken: l.timeTaken,
			Err:       err,
		})
	}
}

func cleanQuery(query string) string {
	result := regexToCleanQueryToDump.ReplaceAllString(query, "")
	result = strings.ReplaceAll(result, "\n", " ")
	return strings.TrimSpace(result)
}

func NewLogInfoExt(ctx context.Context, cf gox.CrossFunction, query string, config *MySQLConfig, callbacks *Callbacks) LogInfo {
	return NewLogInfo(ctx, query, cf.Logger(), config.EnableSqlQueryLogging, config.EnableSqlQueryMetricLogging, callbacks)
}

// NewLogInfo will return a LogInfo object which is initialized with default and other values
func NewLogInfo(ctx context.Context, query string, logger *zap.Logger, enableSqlQueryLogging bool, enableSqlQueryMetricLogging bool, callbacks *Callbacks) LogInfo {
	return LogInfo{
		ctx:                         ctx,
		startTime:                   time.Now().UnixMilli(),
		name:                        util.GetMethodNameName(5),
		query:                       query,
		cleanQuery:                  cleanQuery(query),
		logger:                      logger,
		callbacks:                   callbacks,
		enableSqlQueryLogging:       enableSqlQueryLogging,
		enableSqlQueryMetricLogging: enableSqlQueryMetricLogging,
	}
}

// StartMetricDump should be called if you want to dump the metrics to console - this is mostly used in debugging or
// perf testing
func StartMetricDump(ctx context.Context, config *MySQLConfig) {
	// Start metric dumping - start only once
	startMetricDumpSyncOnce.Do(func() {

		// Dump all metric every 10 sec
		if config.MetricDumpIntervalSec > 0 {
			go metrics.Log(metrics.DefaultRegistry, time.Duration(config.MetricDumpIntervalSec)*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
		} else {
			go metrics.Log(metrics.DefaultRegistry, 10*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
		}

		// Clear all metrics every 10 min and start fresh - this will avoid leak and also will give you fresh stats
		// of last 10 min
		go func() {
			d := 5 * time.Minute
			if config.MetricResetAfterEveryNSec > 0 {
				d = time.Duration(config.MetricResetAfterEveryNSec) * time.Second
			}
		exit:
			for {
				select {
				case <-ctx.Done():
					goto exit
				case <-time.After(d):
					metrics.DefaultRegistry.UnregisterAll()
				}
			}
		}()
	})
}
