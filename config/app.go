package config

type App struct {
	AppName                     string `json:"name" yaml:"name"`
	HttpPort                    int    `json:"http_port" yaml:"http_port"`
	Environment                 string `json:"env" yaml:"env"`
	RequestReadTimeoutMs        int    `json:"request_read_timeout_ms" yaml:"request_read_timeout_ms"`
	RequestWriteTimeoutMs       int    `json:"request_write_timeout_ms" yaml:"request_write_timeout_ms"`
	OutstandingRequestTimeoutMs int    `json:"outstanding_request_timeout_ms" yaml:"outstanding_request_timeout_ms"`
	IdleTimeoutMs               int    `json:"idle_timeout_ms" yaml:"idle_timeout_ms"`
}

func (a *App) SetupDefaults() {
	if a.RequestReadTimeoutMs == 0 {
		a.RequestReadTimeoutMs = 1000
	}
	if a.RequestWriteTimeoutMs == 0 {
		a.RequestWriteTimeoutMs = 1000
	}
	if a.OutstandingRequestTimeoutMs == 0 {
		a.OutstandingRequestTimeoutMs = 5000
	}
	if a.IdleTimeoutMs == 0 {
		a.IdleTimeoutMs = 1000
	}
}
