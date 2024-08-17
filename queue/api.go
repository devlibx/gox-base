package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlibx/gox-base/v2"
	errors2 "github.com/devlibx/gox-base/v2/errors"
	"github.com/oklog/ulid/v2"
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

	SubStatusInternalError           = StatusFailed*10 + 1
	SubStatusApplicationError        = StatusFailed*10 + 2
	SubStatusNoRetryPendingError     = StatusFailed*10 + 3
	SubStatusRetryPendingError       = StatusFailed*10 + 4
	SubStatusRetryIgnoredByUserError = StatusFailed*10 + 5
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

	// Schedule method puts a request on the queue to be executed at a scheduled time.
	// It takes a context and a ScheduleRequest as input and returns a ScheduleResponse or an error.
	Schedule(ctx context.Context, req ScheduleRequest) (*ScheduleResponse, error)

	// Poll method retrieves a request from the queue to be executed immediately.
	// It takes a context and a PollRequest as input and returns a PollResponse or an error.
	Poll(ctx context.Context, req PollRequest) (*PollResponse, error)

	// FetchJobDetails retrieves details of a specific job from the queue.
	// It takes a context and a JobDetailsRequest as input and returns a JobDetailsResponse or an error.
	FetchJobDetails(ctx context.Context, req JobDetailsRequest) (result *JobDetailsResponse, err error)

	// MarkJobFailedAndScheduleRetry marks a job as failed and schedules it for retry.
	// It takes a context and a MarkJobFailedWithRetryRequest as input and returns a MarkJobFailedWithRetryResponse or an error.
	MarkJobFailedAndScheduleRetry(ctx context.Context, req MarkJobFailedWithRetryRequest) (result *MarkJobFailedWithRetryResponse, err error)

	// MarkJobCompleted marks a job as completed.
	// It takes a context and a MarkJobCompletedRequest as input and returns a MarkJobCompletedResponse or an error.
	MarkJobCompleted(ctx context.Context, req MarkJobCompletedRequest) (result *MarkJobCompletedResponse, err error)

	// UpdateJobData updates the data for the given job
	// It takes a context and a UpdateJobDataRequest as input and returns a UpdateJobDataResponse or an error.
	UpdateJobData(ctx context.Context, req UpdateJobDataRequest) (result *UpdateJobDataResponse, err error)
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

	InternalTx           *sql.Tx
	InternalRetryGroupId string
}

func (s ScheduleRequest) String() string {
	return fmt.Sprintf(
		"ScheduleRequest{At:%s, JobType:%d, Tenant:%d, CorrelationId:%s, RemainingExecution:%d, StringUdf1:%s, StringUdf2:%s, IntUdf1:%d, IntUdf2:%d, Properties:%v}",
		s.At, s.JobType, s.Tenant, s.CorrelationId, s.RemainingExecution, s.StringUdf1, s.StringUdf2, s.IntUdf1, s.IntUdf2, s.Properties,
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

type PollResponseError struct {
	WaitForDurationBeforeTrying       time.Duration
	NextJobTimeAvailableForProcessing time.Time
}

func (p PollResponseError) Error() string {
	return fmt.Sprintf("(PollResponseError) no job avaliable to process now. Wait for %vms", p.WaitForDurationBeforeTrying.Milliseconds())
}

// JobDetailsRequest response of schedule
type JobDetailsRequest struct {
	Id string
}

// MarkJobFailedWithRetryRequest mark it failed and set for retry
type MarkJobFailedWithRetryRequest struct {
	Id              string
	ScheduleRetryAt time.Time
}

// MarkJobFailedWithRetryResponse mark it failed and set for retry
type MarkJobFailedWithRetryResponse struct {
	RetryJobId string
	Done       bool
}

type MarkJobCompletedRequest struct {
	Id string
}

type MarkJobCompletedResponse struct {
}

// JobDetailsResponse response of schedule
type JobDetailsResponse struct {
	Id string

	At time.Time

	// Job types
	JobType  int
	State    int
	SubState int

	// Tenant - default is 0
	Tenant int

	// CorrelationId will help jobs to be linked together - when a job succeeds it will mark all jobs to be completed
	CorrelationId string
	RetryGroup    string

	// How many times this job can run (1 time for each + max no of retries)
	// e.g. If it is set 4 then it will run once and in case of error it will be retried 3 times
	RemainingExecution int

	// UDF for application usage - these will be indexed
	StringUdf1 string
	StringUdf2 string
	IntUdf1    int
	IntUdf2    int

	Properties map[string]interface{}
}

func (s PollResponse) String() string {
	return fmt.Sprintf("PollResponse{Id:%s, ProcessAtTimeUsed:%s, RecordPartitionTime:%s}", s.Id, s.ProcessAtTimeUsed.Local().Format(time.RFC3339), s.RecordPartitionTime.Local().Format(time.RFC3339))
}

type UpdateJobDataRequest struct {
	Id         string
	StringUdf1 string
	StringUdf2 string
	IntUdf1    int
	IntUdf2    int
	Properties map[string]interface{}
}

type UpdateJobDataResponse struct {
}

// MySqlBackedStoreBackendConfig is the config to be used for MySQL backed queue
type MySqlBackedStoreBackendConfig struct {
	Host                 string              `json:"host,omitempty"`
	Port                 int                 `json:"port,omitempty"`
	User                 string              `json:"user,omitempty"`
	Password             string              `json:"password,omitempty"`
	Database             string              `json:"database,omitempty"`
	MaxIdleConnection    int                 `json:"max_idle_connection"`
	MaxOpenConnection    int                 `json:"max_open_connection"`
	ConnMaxLifetimeInSec int                 `json:"conn_max_lifetime_in_sec"`
	Properties           gox.StringObjectMap `json:"properties,omitempty"`
	ColumnMapping        map[string]string   `json:"column_mapping,omitempty"`
}

func (m *MySqlBackedStoreBackendConfig) SetupDefault() {
	if m.Host == "" {
		m.Host = "localhost"
	}
	if m.Port <= 0 {
		m.Port = 3306
	}
	if m.ColumnMapping == nil {
		m.ColumnMapping = map[string]string{}
	}
	if m.MaxIdleConnection <= 0 {
		m.MaxIdleConnection = 10
	}
	if m.MaxOpenConnection <= 0 {
		m.MaxOpenConnection = 10
	}
	if m.ConnMaxLifetimeInSec <= 0 {
		m.ConnMaxLifetimeInSec = 60
	}
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

func GeneratePartitionTimeByRecordId(id string) (time.Time, error) {
	if t, err := RecordIdToTime(id); err != nil {
		return time.Time{}, err
	} else {
		return InternalImplEndOfWeek(t), nil
	}
}

func RecordIdToTime(id string) (time.Time, error) {
	if i, err := ulid.Parse(id); err == nil {
		return time.UnixMilli(int64(i.Time())), nil
	} else {
		return time.Time{}, errors2.Wrap(err, "failed to get time from id: id=%s", i)
	}
}

func InternalImplEndOfWeek(inputTime time.Time) time.Time {
	if true {
		return inputTime
	}
	inputTime = inputTime.Truncate(time.Hour).Add(time.Duration(-1 * inputTime.Hour()))
	daysUntilSunday := 0
	switch inputTime.Weekday() {
	case time.Monday:
		daysUntilSunday = 6
		break
	case time.Tuesday:
		daysUntilSunday = 5
		break
	case time.Wednesday:
		daysUntilSunday = 4
		break
	case time.Thursday:
		daysUntilSunday = 3
		break
	case time.Friday:
		daysUntilSunday = 2
		break
	case time.Saturday:
		daysUntilSunday = 1
		break
	case time.Sunday:
		daysUntilSunday = 0
		break
	}

	// Use the Add method to add the remaining days to the input time.

	endOfWeekTime := inputTime.Add(time.Duration(daysUntilSunday) * 24 * time.Hour)

	// Set the time to the end of the day (23:59:59).
	endOfWeekTime = time.Date(endOfWeekTime.Year(), endOfWeekTime.Month(), endOfWeekTime.Day(), 23, 59, 59, 0, time.Now().Location())

	return endOfWeekTime
}
