package config

import (
	"fmt"
	"github.com/devlibx/gox-base/serialization"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// ServerConfig
type ServerConfig struct {
	Apis    Apis    `yaml:"apis"`
	Servers Servers `yaml:"servers"`
}

// Apis
type Apis struct {
	GetPost GetPost `yaml:"getPost"`
}

// GetPost
type GetPost struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Timeout int    `yaml:"timeout"`
}

// Servers
type Servers struct {
	TestServer TestServer `yaml:"testServer"`
}

// TestServer
type TestServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Yaml2Go
type Yaml2Go struct {
	Client       Client       `yaml:"client"`
	ServerConfig ServerConfig `yaml:"server_config"`
}

// Client
type Client struct {
	Enabled bool    `yaml:"enabled"`
	Id      string  `yaml:"id"`
	Price   float64 `yaml:"price"`
}

func TestReadParameterizedConfigYamlFile_Prod(t *testing.T) {

	_ = os.Setenv("PRDO_testServer", "test.prod")
	_ = os.Setenv("STAGE_testServer", "test.stage")

	yaml2Go := Yaml2Go{}
	err := ReadParameterizedYamlFile("../testdata/app_with_env_var_and_params.yml", &yaml2Go, "prod")
	assert.NoError(t, err)

	assert.Equal(t, "prod_client", yaml2Go.Client.Id)
	assert.Equal(t, 10.001, yaml2Go.Client.Price)
	assert.Equal(t, true, yaml2Go.Client.Enabled)

	assert.Equal(t, "POST", yaml2Go.ServerConfig.Apis.GetPost.Method)
	assert.Equal(t, "/get", yaml2Go.ServerConfig.Apis.GetPost.Path)
	assert.Equal(t, 10, yaml2Go.ServerConfig.Apis.GetPost.Timeout)

	assert.Equal(t, "test.prod", yaml2Go.ServerConfig.Servers.TestServer.Host)
	assert.Equal(t, 80, yaml2Go.ServerConfig.Servers.TestServer.Port)

	yaml, err := serialization.ToYaml(yaml2Go)
	assert.NoError(t, err)
	fmt.Println(yaml)
}

func TestReadParameterizedConfigYamlFile_Stage(t *testing.T) {

	_ = os.Setenv("PRDO_testServer", "test.prod")
	_ = os.Setenv("STAGE_testServer", "test.stage")

	yaml2Go := Yaml2Go{}
	err := ReadParameterizedYamlFile("../testdata/app_with_env_var_and_params.yml", &yaml2Go, "stage")
	assert.NoError(t, err)

	assert.Equal(t, "stage_client", yaml2Go.Client.Id)
	assert.Equal(t, 10.002, yaml2Go.Client.Price)
	assert.Equal(t, false, yaml2Go.Client.Enabled)

	assert.Equal(t, "POST", yaml2Go.ServerConfig.Apis.GetPost.Method)
	assert.Equal(t, "/get", yaml2Go.ServerConfig.Apis.GetPost.Path)
	assert.Equal(t, 20, yaml2Go.ServerConfig.Apis.GetPost.Timeout)

	assert.Equal(t, "test.stage", yaml2Go.ServerConfig.Servers.TestServer.Host)
	assert.Equal(t, 80, yaml2Go.ServerConfig.Servers.TestServer.Port)
}

func TestReadParameterizedConfigYamlFile_Default(t *testing.T) {

	_ = os.Setenv("PRDO_testServer", "test.prod")
	_ = os.Setenv("STAGE_testServer", "test.stage")

	yaml2Go := Yaml2Go{}
	err := ReadParameterizedYamlFile("../testdata/app_with_env_var_and_params.yml", &yaml2Go, "random")
	assert.NoError(t, err)

	assert.Equal(t, "random_client", yaml2Go.Client.Id)
	assert.Equal(t, 10.004, yaml2Go.Client.Price)
	assert.Equal(t, false, yaml2Go.Client.Enabled)

	assert.Equal(t, "POST", yaml2Go.ServerConfig.Apis.GetPost.Method)
	assert.Equal(t, "/get", yaml2Go.ServerConfig.Apis.GetPost.Path)
	assert.Equal(t, 1000, yaml2Go.ServerConfig.Apis.GetPost.Timeout)

	assert.Equal(t, "localhost", yaml2Go.ServerConfig.Servers.TestServer.Host)
	assert.Equal(t, 8090, yaml2Go.ServerConfig.Servers.TestServer.Port)
}
