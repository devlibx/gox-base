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

var WorkflowActionWorkflowStarted       = workflowAction("add-workflow-started-event")

// Pre-defined values for TagSysComponent
var (
	ComponentTaskQueue                = component("taskqueue")
	ComponentHistoryEngine            = component("history-engine")
	ComponentHistoryCache             = component("history-cache")
	ComponentEventsCache              = component("events-cache")
	ComponentTransferQueue            = component("transfer-queue-processor")
	ComponentVisibilityQueue          = component("visibility-queue-processor")
	ComponentTimerQueue               = component("timer-queue-processor")
	ComponentTimerBuilder             = component("timer-builder")
	ComponentReplicatorQueue          = component("replicator-queue-processor")
	ComponentShardController          = component("shard-controller")
	ComponentShard                    = component("shard")
	ComponentShardItem                = component("shard-item")
	ComponentShardEngine              = component("shard-engine")
	ComponentMatchingEngine           = component("matching-engine")
	ComponentReplicator               = component("replicator")
	ComponentReplicationTaskProcessor = component("replication-task-processor")
	ComponentHistoryReplicator        = component("history-replicator")
	ComponentIndexer                  = component("indexer")
	ComponentIndexerProcessor         = component("indexer-processor")
	ComponentIndexerESProcessor       = component("indexer-es-processor")
	ComponentESVisibilityManager      = component("es-visibility-manager")
	ComponentArchiver                 = component("archiver")
	ComponentBatcher                  = component("batcher")
	ComponentWorker                   = component("worker")
	ComponentServiceResolver          = component("service-resolver")
	ComponentMetadataInitializer      = component("metadata-initializer")
	VersionChecker                    = component("version-checker")
)

// Pre-defined values for TagSysLifecycle
var (
	LifeCycleStarting         = lifecycle("Starting")
	LifeCycleStarted          = lifecycle("Started")
	LifeCycleStopping         = lifecycle("Stopping")
	LifeCycleStopped          = lifecycle("Stopped")
	LifeCycleStopTimedout     = lifecycle("StopTimedout")
	LifeCycleStartFailed      = lifecycle("StartFailed")
	LifeCycleStopFailed       = lifecycle("StopFailed")
	LifeCycleProcessingFailed = lifecycle("ProcessingFailed")
)

// Pre-defined values for SysErrorType
var (
	ErrorTypeInvalidHistoryAction         = errorType("InvalidHistoryAction")
	ErrorTypeInvalidQueryTask             = errorType("InvalidQueryTask")
	ErrorTypeQueryTaskFailed              = errorType("QueryTaskFailed")
	ErrorTypePersistentStoreError         = errorType("PersistentStoreError")
	ErrorTypeHistorySerializationError    = errorType("HistorySerializationError")
	ErrorTypeHistoryDeserializationError  = errorType("HistoryDeserializationError")
	ErrorTypeDuplicateTask                = errorType("DuplicateTask")
	ErrorTypeMultipleCompletionCommands   = errorType("MultipleCompletionCommands")
	ErrorTypeDuplicateTransferTask        = errorType("DuplicateTransferTask")
	ErrorTypeWorkflowTaskFailed           = errorType("WorkflowTaskFailed")
	ErrorTypeInvalidMutableStateAction    = errorType("InvalidMutableStateAction")
	ErrorTypeInvalidMemWorkflowTaskAction = errorType("InvalidMemWorkflowTaskAction")
)

// Pre-defined values for SysShardUpdate
var (
	// Shard context events
	ValueShardRangeUpdated            = shardupdate("ShardRangeUpdated")
	ValueShardAllocateTimerBeforeRead = shardupdate("ShardAllocateTimerBeforeRead")
	ValueRingMembershipChangedEvent   = shardupdate("RingMembershipChangedEvent")
)

// Pre-defined values for OperationResult
var (
	OperationFailed   = operationResult("OperationFailed")
	OperationStuck    = operationResult("OperationStuck")
	OperationCritical = operationResult("OperationCritical")
)

// Pre-defined values for TagSysStoreOperation
var (
	StoreOperationGetTasks                = storeOperation("get-tasks")
	StoreOperationCompleteTask            = storeOperation("complete-task")
	StoreOperationCompleteTasksLessThan   = storeOperation("complete-tasks-less-than")
	StoreOperationCreateWorkflowExecution = storeOperation("create-wf-execution")
	StoreOperationGetWorkflowExecution    = storeOperation("get-wf-execution")
	StoreOperationUpdateWorkflowExecution = storeOperation("update-wf-execution")
	StoreOperationDeleteWorkflowExecution = storeOperation("delete-wf-execution")
	StoreOperationUpdateShard             = storeOperation("update-shard")
	StoreOperationCreateTask              = storeOperation("create-task")
	StoreOperationUpdateTaskQueue         = storeOperation("update-task-queue")
	StoreOperationStopTaskQueue           = storeOperation("stop-task-queue")
)
