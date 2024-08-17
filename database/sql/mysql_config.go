package goxSql

import (
	"context"
	"github.com/devlibx/gox-base/v2/util"
)

const DbCallNameKeyInCyx = "__SQLCX_DB_CALL_NAME__"

type MySQLConfig struct {
	ServerName   string `json:"server_name" yaml:"server_name"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	Db           string `json:"db" yaml:"db"`
	TablePrefix  string
	TablePostfix string

	EnableSqlQueryLogging bool `json:"enable_sql_query_logging" yaml:"enable_sql_query_logging"`

	EnableSqlQueryMetricLogging bool `json:"enable_sql_query_metric_logging" yaml:"enable_sql_query_metric_logging"`
	MetricDumpIntervalSec       int  `json:"metric_dump_interval_sec" yaml:"metric_dump_interval_sec"`
	MetricResetAfterEveryNSec   int  `json:"metric_reset_after_every_n_sec" yaml:"metric_reset_after_every_n_sec"`
}

func (m *MySQLConfig) SetupDefaults() {
	if util.IsStringEmpty(m.Host) {
		m.Host = "localhost"
	}
	if m.Port <= 0 {
		m.Port = 3306
	}
	if util.IsStringEmpty(m.User) {
		m.User = "test"
	}
	if util.IsStringEmpty(m.Password) {
		m.Password = "test"
	}
	if util.IsStringEmpty(m.Db) {
		m.Db = "conversation"
	}
}

type Callbacks struct {
	PostCallbackFunc PostCallbackFunc
}

type PostCallbackData struct {
	Ctx       context.Context `json:"-"`
	Name      string          `json:"name"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	TimeTaken int64           `json:"time_taken"`
	Err       error           `json:"error"`
}

func (p *PostCallbackData) GetDbCallNameForTracing() string {
	if p.Ctx != nil && p.Ctx.Value(DbCallNameKeyInCyx) != nil {
		if val, ok := p.Ctx.Value(DbCallNameKeyInCyx).(string); ok {
			return "Slow_Query_Trace__" + val
		}
	}
	return "Slow_Query_Trace__" + p.Name
}

type PostCallbackFunc func(data PostCallbackData)
