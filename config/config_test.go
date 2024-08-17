package config

import (
	"github.com/devlibx/gox-base/v2/serialization"
	"github.com/stretchr/testify/assert"
	"testing"
)

type configToTest struct {
	App    App    `json:"app" yaml:"app"`
	Logger Logger `json:"logger" yaml:"logger"`
}

type extendedAppObject struct {
	App  `yaml:",inline"`
	Some string `json:"some"`
}

type configToTestWithExtendedAppObject struct {
	App    extendedAppObject `json:"app" yaml:"app"`
	Logger Logger            `json:"logger" yaml:"logger"`
}

func TestAppConfig_Yaml(t *testing.T) {
	// Try with normal struct
	conf := &configToTest{}
	err := serialization.ReadYamlFromString(yamlConfig, conf)
	assert.NoError(t, err)

	validateApp(conf, t)
	validateLogger(conf, t)

	// Try with modified struct
	confSub := &configToTestWithExtendedAppObject{}
	err = serialization.ReadYamlFromString(yamlConfig, confSub)
	assert.NoError(t, err)

	validateAppSubConf(confSub, t)
	validateLoggerSubConf(confSub, t)
}

func TestAppConfig_Json(t *testing.T) {
	// Try with normal struct
	conf := &configToTest{}
	err := serialization.JsonToObject(jsonString, conf)
	assert.NoError(t, err)

	validateApp(conf, t)
	validateLogger(conf, t)

	// Try with modified struct
	confSub := &configToTestWithExtendedAppObject{}
	err = serialization.JsonToObject(jsonString, confSub)
	assert.NoError(t, err)

	validateAppSubConf(confSub, t)
	validateLoggerSubConf(confSub, t)

}

func validateApp(conf *configToTest, t *testing.T) {
	assert.Equal(t, "test", conf.App.AppName)
	assert.Equal(t, "dev", conf.App.Environment)
	assert.Equal(t, 8080, conf.App.HttpPort)
}

func validateLogger(conf *configToTest, t *testing.T) {
	assert.Equal(t, "debug", conf.Logger.LogLevel)
	assert.True(t, conf.Logger.EnableConsoleLog)
}

func validateAppSubConf(conf *configToTestWithExtendedAppObject, t *testing.T) {
	assert.Equal(t, "test", conf.App.AppName)
	assert.Equal(t, "dev", conf.App.Environment)
	assert.Equal(t, 8080, conf.App.HttpPort)
	assert.Equal(t, "dummy", conf.App.Some)
}

func validateLoggerSubConf(conf *configToTestWithExtendedAppObject, t *testing.T) {
	assert.Equal(t, "debug", conf.Logger.LogLevel)
	assert.True(t, conf.Logger.EnableConsoleLog)
}

var yamlConfig = `
app:
  name: test
  http_port: 8080
  env: dev
  some: dummy
logger:
  level: debug
  enable_console_log: true
`

var jsonString = `
{
	"app": {
		"name": "test",
		"http_port": 8080,
		"env": "dev",
		"some": "dummy"
	},
	"logger": {
		"level": "debug",
		"enable_console_log": true
	}
}
`
