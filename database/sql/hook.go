package goxSql

import (
	"context"
	"github.com/devlibx/gox-base/v2/util"
	"time"
)

// LogInfoFunc is used to generate LogInfo for defer
type LogInfoFunc func(ctx context.Context, query string, args ...interface{}) LogInfo

// DefaultLogInfoFunc is a default NoOp method
var DefaultLogInfoFunc LogInfoFunc = func(ctx context.Context, query string, args ...interface{}) LogInfo {
	return LogInfo{
		ctx:                         ctx,
		startTime:                   time.Now().UnixMilli(),
		name:                        util.GetMethodNameName(5),
		query:                       query,
		cleanQuery:                  cleanQuery(query),
		callbacks:                   nil,
		enableSqlQueryLogging:       true,
		enableSqlQueryMetricLogging: true,
	}
}

// RegisterLogInfoFunc will help you to register your function
func RegisterLogInfoFunc(f LogInfoFunc) {
	DefaultLogInfoFunc = f
}

func BuildNewLogInfo(ctx context.Context, query string, args ...interface{}) LogInfo {
	if DefaultLogInfoFunc == nil {
		return LogInfo{}
	} else {
		return DefaultLogInfoFunc(ctx, query, args...)
	}
}
