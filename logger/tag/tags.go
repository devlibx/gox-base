// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tag

import (
	"fmt"
	"time"
)

// All logging tags are defined in this file.
// To help finding available tags, we recommend that all tags to be categorized and placed in the corresponding section.
// We currently have those categories:
//   0. Common tags that can't be categorized(or belong to more than one)
//   1. Workflow: these tags are information that are useful to our customer, like workflow-id/run-id/task-queue/...
//   2. System : these tags are internal information which usually cannot be understood by our customers,

// LoggingCallAtKey is reserved tag
const LoggingCallAtKey = "logging-call-at"

///////////////////  Common tags defined here ///////////////////

// Operation returns tag for Operation
func Operation(operation string) ZapTag {
	return NewStringTag("operation", operation)
}

// Error returns tag for Error
func Error(err error) ZapTag {
	return NewErrorTag(err)
}

// ClusterName returns tag for ClusterName
func ClusterName(clusterName string) ZapTag {
	return NewStringTag("cluster-name", clusterName)
}

// Timestamp returns tag for Timestamp
func Timestamp(timestamp time.Time) ZapTag {
	return NewTimeTag("timestamp", timestamp)
}

///////////////////  System tags defined here:  ///////////////////
// Tags with pre-define values

// Component returns tag for Component
func component(component string) ZapTag {
	return NewStringTag("component", component)
}

// Lifecycle returns tag for Lifecycle
func lifecycle(lifecycle string) ZapTag {
	return NewStringTag("lifecycle", lifecycle)
}

// StoreOperation returns tag for StoreOperation
func storeOperation(storeOperation string) ZapTag {
	return NewStringTag("store-operation", storeOperation)
}

// OperationResult returns tag for OperationResult
func operationResult(operationResult string) ZapTag {
	return NewStringTag("operation-result", operationResult)
}

// ErrorType returns tag for ErrorType
func errorType(errorType string) ZapTag {
	return NewStringTag("error-type", errorType)
}

// Shardupdate returns tag for Shardupdate
func shardupdate(shardupdate string) ZapTag {
	return NewStringTag("shard-update", shardupdate)
}

// general

// Service returns tag for Service
func Service(sv string) ZapTag {
	return NewStringTag("service", sv)
}

// Addresses returns tag for Addresses
func Addresses(ads []string) ZapTag {
	return NewStringsTag("addresses", ads)
}

// ListenerName returns tag for ListenerName
func ListenerName(name string) ZapTag {
	return NewStringTag("listener-name", name)
}

// Address return tag for Address
func Address(ad string) ZapTag {
	return NewStringTag("address", ad)
}

// HostID return tag for HostID
func HostID(hid string) ZapTag {
	return NewStringTag("hostId", hid)
}

// Env return tag for runtime environment
func Env(env string) ZapTag {
	return NewStringTag("env", env)
}

// Key returns tag for Key
func Key(k string) ZapTag {
	return NewStringTag("key", k)
}

// Name returns tag for Name
func Name(k string) ZapTag {
	return NewStringTag("name", k)
}

// Value returns tag for Value
func Value(v interface{}) ZapTag {
	return NewAnyTag("value", v)
}

// ValueType returns tag for ValueType
func ValueType(v interface{}) ZapTag {
	return NewStringTag("value-type", fmt.Sprintf("%T", v))
}

// DefaultValue returns tag for DefaultValue
func DefaultValue(v interface{}) ZapTag {
	return NewAnyTag("default-value", v)
}

// IgnoredValue returns tag for IgnoredValue
func IgnoredValue(v interface{}) ZapTag {
	return NewAnyTag("ignored-value", v)
}

// Port returns tag for Port
func Port(p int) ZapTag {
	return NewInt("port", p)
}

// CursorTimestamp returns tag for CursorTimestamp
func CursorTimestamp(timestamp time.Time) ZapTag {
	return NewTimeTag("cursor-timestamp", timestamp)
}

// MetricScope returns tag for MetricScope
func MetricScope(metricScope int) ZapTag {
	return NewInt("metric-scope", metricScope)
}

// StoreType returns tag for StoreType
func StoreType(storeType string) ZapTag {
	return NewStringTag("store-type", storeType)
}

// DetailInfo returns tag for DetailInfo
func DetailInfo(i string) ZapTag {
	return NewStringTag("detail-info", i)
}

// Counter returns tag for Counter
func Counter(c int) ZapTag {
	return NewInt("counter", c)
}

// Number returns tag for Number
func Number(n int64) ZapTag {
	return NewInt64("number", n)
}

// NextNumber returns tag for NextNumber
func NextNumber(n int64) ZapTag {
	return NewInt64("next-number", n)
}

// Bool returns tag for Bool
func Bool(b bool) ZapTag {
	return NewBoolTag("bool", b)
}

// history engine shard

// ShardID returns tag for ShardID
func ShardID(shardID int32) ZapTag {
	return NewInt32("shard-id", shardID)
}

// ShardItem returns tag for ShardItem
func ShardItem(shardItem interface{}) ZapTag {
	return NewAnyTag("shard-item", shardItem)
}

// ShardTime returns tag for ShardTime
func ShardTime(shardTime interface{}) ZapTag {
	return NewAnyTag("shard-time", shardTime)
}

// ShardReplicationAck returns tag for ShardReplicationAck
func ShardReplicationAck(shardReplicationAck int64) ZapTag {
	return NewInt64("shard-replication-ack", shardReplicationAck)
}

// PreviousShardRangeID returns tag for PreviousShardRangeID
func PreviousShardRangeID(id int64) ZapTag {
	return NewInt64("previous-shard-range-id", id)
}

// ShardRangeID returns tag for ShardRangeID
func ShardRangeID(id int64) ZapTag {
	return NewInt64("shard-range-id", id)
}

// ReadLevel returns tag for ReadLevel
func ReadLevel(lv int64) ZapTag {
	return NewInt64("read-level", lv)
}

// MinLevel returns tag for MinLevel
func MinLevel(lv int64) ZapTag {
	return NewInt64("min-level", lv)
}

// MaxLevel returns tag for MaxLevel
func MaxLevel(lv int64) ZapTag {
	return NewInt64("max-level", lv)
}

// ShardTransferAcks returns tag for ShardTransferAcks
func ShardTransferAcks(shardTransferAcks interface{}) ZapTag {
	return NewAnyTag("shard-transfer-acks", shardTransferAcks)
}

// ShardTimerAcks returns tag for ShardTimerAcks
func ShardTimerAcks(shardTimerAcks interface{}) ZapTag {
	return NewAnyTag("shard-timer-acks", shardTimerAcks)
}

// task queue processor

// Task returns tag for Task
func Task(task interface{}) ZapTag {
	return NewAnyTag("queue-task", task)
}

// TaskID returns tag for TaskID
func TaskID(taskID int64) ZapTag {
	return NewInt64("queue-task-id", taskID)
}

// TaskVersion returns tag for TaskVersion
func TaskVersion(taskVersion int64) ZapTag {
	return NewInt64("queue-task-version", taskVersion)
}

// TaskVisibilityTimestamp returns tag for task visibilityTimestamp
func TaskVisibilityTimestamp(timestamp int64) ZapTag {
	return NewInt64("queue-task-visibility-timestamp", timestamp)
}

// NumberProcessed returns tag for NumberProcessed
func NumberProcessed(n int) ZapTag {
	return NewInt("number-processed", n)
}

// NumberDeleted returns tag for NumberDeleted
func NumberDeleted(n int) ZapTag {
	return NewInt("number-deleted", n)
}

// TimerTaskStatus returns tag for TimerTaskStatus
func TimerTaskStatus(timerTaskStatus int32) ZapTag {
	return NewInt32("timer-task-status", timerTaskStatus)
}

// retry

// Attempt returns tag for Attempt
func Attempt(attempt int32) ZapTag {
	return NewInt32("attempt", attempt)
}

// AttemptCount returns tag for AttemptCount
func AttemptCount(attemptCount int64) ZapTag {
	return NewInt64("attempt-count", attemptCount)
}

// AttemptStart returns tag for AttemptStart
func AttemptStart(attemptStart time.Time) ZapTag {
	return NewTimeTag("attempt-start", attemptStart)
}

// AttemptEnd returns tag for AttemptEnd
func AttemptEnd(attemptEnd time.Time) ZapTag {
	return NewTimeTag("attempt-end", attemptEnd)
}

// ScheduleAttempt returns tag for ScheduleAttempt
func ScheduleAttempt(scheduleAttempt int32) ZapTag {
	return NewInt32("schedule-attempt", scheduleAttempt)
}

// ElasticSearch

// ESRequest returns tag for ESRequest
func ESRequest(ESRequest string) ZapTag {
	return NewStringTag("es-request", ESRequest)
}

// ESResponseStatus returns tag for ESResponse status
func ESResponseStatus(status int) ZapTag {
	return NewInt("es-response-status", status)
}

// ESResponseError returns tag for ESResponse error
func ESResponseError(msg string) ZapTag {
	return NewStringTag("es-response-error", msg)
}

// ESKey returns tag for ESKey
func ESKey(ESKey string) ZapTag {
	return NewStringTag("es-mapping-key", ESKey)
}

// ESValue returns tag for ESValue
func ESValue(ESValue []byte) ZapTag {
	// convert value to string type so that the value logged is human readable
	return NewStringTag("es-mapping-value", string(ESValue))
}

// ESConfig returns tag for ESConfig
func ESConfig(c interface{}) ZapTag {
	return NewAnyTag("es-config", c)
}

// ESField returns tag for ESField
func ESField(ESField string) ZapTag {
	return NewStringTag("es-Field", ESField)
}

// ESDocID returns tag for ESDocID
func ESDocID(id string) ZapTag {
	return NewStringTag("es-doc-id", id)
}

// SysStackTrace returns tag for SysStackTrace
func SysStackTrace(stackTrace string) ZapTag {
	return NewStringTag("sys-stack-trace", stackTrace)
}

// TokenLastEventID returns tag for TokenLastEventID
func TokenLastEventID(id int64) ZapTag {
	return NewInt64("token-last-event-id", id)
}

///////////////////  XDC tags defined here: xdc- ///////////////////

// SourceCluster returns tag for SourceCluster
func SourceCluster(sourceCluster string) ZapTag {
	return NewStringTag("xdc-source-cluster", sourceCluster)
}

// PrevActiveCluster returns tag for PrevActiveCluster
func PrevActiveCluster(prevActiveCluster string) ZapTag {
	return NewStringTag("xdc-prev-active-cluster", prevActiveCluster)
}

// FailoverMsg returns tag for FailoverMsg
func FailoverMsg(failoverMsg string) ZapTag {
	return NewStringTag("xdc-failover-msg", failoverMsg)
}

// FailoverVersion returns tag for Version
func FailoverVersion(version int64) ZapTag {
	return NewInt64("xdc-failover-version", version)
}

// CurrentVersion returns tag for CurrentVersion
func CurrentVersion(currentVersion int64) ZapTag {
	return NewInt64("xdc-current-version", currentVersion)
}

// IncomingVersion returns tag for IncomingVersion
func IncomingVersion(incomingVersion int64) ZapTag {
	return NewInt64("xdc-incoming-version", incomingVersion)
}

// FirstEventVersion returns tag for FirstEventVersion
func FirstEventVersion(version int64) ZapTag {
	return NewInt64("xdc-first-event-version", version)
}

// LastEventVersion returns tag for LastEventVersion
func LastEventVersion(version int64) ZapTag {
	return NewInt64("xdc-last-event-version", version)
}

// TokenLastEventVersion returns tag for TokenLastEventVersion
func TokenLastEventVersion(version int64) ZapTag {
	return NewInt64("xdc-token-last-event-version", version)
}

// TransportType returns tag for transportType
func TransportType(transportType string) ZapTag {
	return NewStringTag("transport-type", transportType)
}

// ActivityInfo returns tag for activity info
func ActivityInfo(activityInfo interface{}) ZapTag {
	return NewAnyTag("activity-info", activityInfo)
}

// WorkflowTaskRequestId returns tag for workflow task RequestId
func WorkflowTaskRequestId(s string) ZapTag {
	return NewStringTag("workflow-task-request-id", s)
}

// AckLevel returns tag for ack level
func AckLevel(s interface{}) ZapTag {
	return NewAnyTag("ack-level", s)
}

// MinQueryLevel returns tag for query level
func MinQueryLevel(s time.Time) ZapTag {
	return NewTimeTag("min-query-level", s)
}

// MaxQueryLevel returns tag for query level
func MaxQueryLevel(s time.Time) ZapTag {
	return NewTimeTag("max-query-level", s)
}

// BootstrapHostPorts returns tag for bootstrap host ports
func BootstrapHostPorts(s string) ZapTag {
	return NewStringTag("bootstrap-hostports", s)
}

// TLSCertFile returns tag for TLS cert file name
func TLSCertFile(filePath string) ZapTag {
	return NewStringTag("tls-cert-file", filePath)
}

// TLSKeyFile returns tag for TLS key file
func TLSKeyFile(filePath string) ZapTag {
	return NewStringTag("tls-key-file", filePath)
}

// TLSCertFiles returns tag for TLS cert file names
func TLSCertFiles(filePaths []string) ZapTag {
	return NewStringsTag("tls-cert-files", filePaths)
}

// WorkflowScheduleID returns tag for WorkflowScheduleID
func WorkflowScheduleID(scheduleID int64) ZapTag {
	return NewInt64("wf-schedule-id", scheduleID)
}

// WorkflowNamespace returns tag for WorkflowNamespace
func WorkflowNamespace(namespace string) ZapTag {
	return NewStringTag("wf-namespace", namespace)
}

// WorkflowAction returns tag for WorkflowAction
func workflowAction(action string) ZapTag {
	return NewStringTag("wf-action", action)
}
