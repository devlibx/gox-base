package config

type App struct {
	AppName     string `json:"name" yaml:"name"`
	HttpPort    int    `json:"http_port" yaml:"http_port"`
	Environment string `json:"env" yaml:"env"`
}
