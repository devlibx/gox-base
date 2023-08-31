package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlibx/gox-base"
	"strings"
	"time"
)

// All status to be used for job status
var (
	StatusScheduled  = 1
	StatusProcessing = 2
	StatusDone       = 3
	StatusFailed     = 4

	SubStatusScheduledOk            = StatusScheduled*10 + 0
	SubStatusDone                   = StatusDone*10 + 0
	SubStatusDoneDueToCorrelatedJob = StatusDone*10 + 1
	SubStatusInternalError          = StatusFailed*10 + 1
	SubStatusNoRetryPendingError    = StatusFailed*10 + 2
)

// ErrNoMoreRetry indicate that no more retries are needed
var ErrNoMoreRetry = errors.New("do not retry anymore")

var NoJobsToRunAtCurrently = errors.New("queue does not have a job to run now")

type MySqlBackedQueueConfig struct {
	Tenant     int `json:"tenant,omitempty"`
	MaxJobType int `json:"max_job_type,omitempty"`

	UsePreparedStatement       bool `json:"use_prepared_statement"`
	UseMinQueryToPickLatestRow bool `json:"use_min_query_to_pick_latest_row"`

	DontRunPoller bool
}

// Queue is an interface to provide all queue related methods. It allows you to schedule, poll etc
type Queue interface {

	// Schedule method will put this request on the queue to be executed on the time
	Schedule(ctx context.Context, req ScheduleRequest) (*ScheduleResponse, error)

	// Poll method will put this request on the queue to be executed on the time
	Poll(ctx context.Context, req PollRequest) (*PollResponse, error)
}

// ScheduleRequest is a request to schedule a run of this job
type ScheduleRequest struct {
	At time.Time

	// Job types
	JobType int

	// Tenant - default is 0
	Tenant int

	// CorrelationId will help jobs to be linked together - when a job succeeds it will mark all jobs to be completed
	CorrelationId string

	// How many times this job can run (1 time for each + max no of retries)
	// e.g. If it is set 4 then it will run once and in case of error it will be retried 3 times
	RemainingExecution int
	RetryBackoffAlgo   RetryBackoffAlgo

	// UDF for application usage - these will be indexed
	StringUdf1 string
	StringUdf2 string
	IntUdf1    int
	IntUdf2    int

	Properties map[string]interface{}
}

func (s ScheduleRequest) String() string {
	return fmt.Sprintf(
		"ScheduleRequest{At:%s, JobType:%s, Tenant:%s, CorrelationId:%s, RemainingExecution:%d, RetryBackoffAlgo:%v, StringUdf1:%s, StringUdf2:%s, IntUdf1:%d, IntUdf2:%d, Properties:%v}",
		s.At, s.JobType, s.Tenant, s.CorrelationId, s.RemainingExecution, s.RetryBackoffAlgo, s.StringUdf1, s.StringUdf2, s.IntUdf1, s.IntUdf2, s.Properties,
	)
}

// ScheduleResponse response of schedule
type ScheduleResponse struct {
	Id string
}

// PollRequest response of schedule
type PollRequest struct {
	Tenant  int
	JobType int
}

// PollResponse response of schedule
type PollResponse struct {
	Id                  string
	RecordPartitionTime time.Time
	ProcessAtTimeUsed   time.Time
}

func (s PollResponse) String() string {
	return fmt.Sprintf("PollResponse{Id:%s, ProcessAtTimeUsed:%s, RecordPartitionTime:%s}", s.Id, s.ProcessAtTimeUsed.Local().Format(time.RFC3339), s.RecordPartitionTime.Local().Format(time.RFC3339))
}

// MySqlBackedStoreBackendConfig is the config to be used for MySQL backed queue
type MySqlBackedStoreBackendConfig struct {
	Host          string              `json:"host,omitempty"`
	Port          int                 `json:"port,omitempty"`
	User          string              `json:"user,omitempty"`
	Password      string              `json:"password,omitempty"`
	Database      string              `json:"database,omitempty"`
	MaxConnection int                 `json:"max_connection,omitempty"`
	MinConnection int                 `json:"min_connection,omitempty"`
	Properties    gox.StringObjectMap `json:"properties,omitempty"`
	ColumnMapping map[string]string   `json:"column_mapping,omitempty"`
}

// StoreBackend is the backend to be used to give connections to store
type StoreBackend interface {

	// Init has to be called at the beginning
	Init() error

	// GetSqlDb the SQL Db to be used for operation
	GetSqlDb() (*sql.DB, error)

	// Close to be called at the end
	Close() error
}

type QueryRewriter interface {
	RewriteQuery(table string, input string) string
}

func NewUdfAndTableNameQueryRewriter(tableName string) QueryRewriter {
	return &UdfAndTableNameQueryRewriter{tableName: tableName}
}

type UdfAndTableNameQueryRewriter struct {
	tableName  string
	udfString1 string
	udfString2 string
	udfInt1    string
	udfInt2    string
}

func (n *UdfAndTableNameQueryRewriter) RewriteQuery(table string, input string) string {
	switch table {
	case "jobs":
		input = strings.ReplaceAll(input, "jobs", n.tableName)
		break

	case "jobs_data":
		input = strings.ReplaceAll(input, "jobs_data", n.tableName+"_data")
		if n.udfString1 != "" {
			input = strings.ReplaceAll(input, "string_udf_1", n.udfString1)
		}
		if n.udfString2 != "" {
			input = strings.ReplaceAll(input, "string_udf_2", n.udfString2)
		}
		if n.udfInt1 != "" {
			input = strings.ReplaceAll(input, "int_udf_1", n.udfInt1)
		}
		if n.udfInt2 != "" {
			input = strings.ReplaceAll(input, "int_udf_2", n.udfInt2)
		}
		break
	}

	return input
}
