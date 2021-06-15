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

var strTestReadParameterizedConfigYamlFile_Prod_List = `
client:
  id: "env:string: prod=prod_client; stage=stage_client; dev=dev_client; default=random_client"
  enabled: "env:bool: prod=true; stage=false; dev=false; default=false"
  price: "env:float: prod=10.001; stage=10.002; dev=10.003; default=10.004"
  option:
    - "env:string: prod=call_prod; stage=call_stage; dev=call_dev; default=call_prod_default"
    - "env:string: prod=sms_prod; stage=sms_stage; dev=sms_dev; default=sms_prod_default"
    - email
    - sub_options:
        key: "env:string: prod=key_prod; default=key_default"
    - true
    - false
    - 10
    - 10.011
`

var testStringTestReadParameterizedConfigYamlListOfMaps = `
client:
  option:
    - call:
        name: "env:string: prod=call_prod; default=call_default"
    - sms:
        name: "env:string: prod=sms_prod; default=sms_default"
    - push:
        name: "env:string: prod=push_prod; default=push_default"    
`

var testStringTestReadParameterizedConfigYamlWithListAndListMap = `
client:
  id: "env:string: prod=$PROD_ID_acbdefgh; stage=$STAGE_ID_acbdefgh; dev=$DEV_ID_acbdefgh; default=$RANDOM_CLIENT_acbdefgh"
  enabled: "env:bool: prod=true; stage=false; dev=false; default=false"
  price: "env:float: prod=10.001; stage=10.002; dev=10.003; default=10.004"
  option:
    - "env:string: prod=call_prod; stage=call_stage; dev=call_dev; default=call_prod_default"
    - "env:string: prod=sms_prod; stage=sms_stage; dev=sms_dev; default=sms_prod_default"
    - email
    - sub_options:
        key: "env:string: prod=key_prod; default=key_default"
    - true
    - false
    - 10
    - 10.011
  option_map:
    - call:
        name: "env:string: prod=call_prod; default=call_default"
    - sms:
        name: "env:string: prod=sms_prod; default=sms_default"
    - push:
        name: "env:string: prod=push_prod; default=push_default"  
`

func TestReadParameterizedConfigYamlWithListAndListMap(t *testing.T) {
	// You will get a env variable from outside but setting it for test here
	_ = os.Setenv("PROD_ID_acbdefgh", "prod_client")

	// ---------------------------------- Wrapper Json Struct ----------------------------------------------------------
	// Internal nodes
	type internalNode struct {
		Name string `yaml:"name"`
	}

	type listMap map[string]internalNode

	// List of options
	type internalList struct {
		Enabled   bool          `yaml:"enabled"`
		Id        string        `yaml:"id"`
		Price     float64       `yaml:"price"`
		Option    []interface{} `yaml:"option"`
		OptionMap []listMap     `yaml:"option_map"`
	}

	type wrapper struct {
		Client internalList `yaml:"client"`
	}

	// Read data in struct for env=prod
	yaml2Go := wrapper{}
	err := ReadParameterizedYaml(testStringTestReadParameterizedConfigYamlWithListAndListMap, &yaml2Go, "prod")
	assert.NoError(t, err)

	// Verify all data
	assert.Equal(t, "prod_client", yaml2Go.Client.Id)
	assert.Equal(t, 10.001, yaml2Go.Client.Price)
	assert.Equal(t, true, yaml2Go.Client.Enabled)

	// Verify all data - option list
	assert.Equal(t, 8, len(yaml2Go.Client.Option))
	assert.Equal(t, "call_prod", yaml2Go.Client.Option[0])
	assert.Equal(t, "sms_prod", yaml2Go.Client.Option[1])
	assert.Equal(t, "email", yaml2Go.Client.Option[2])
	assert.Equal(t, true, yaml2Go.Client.Option[4])
	assert.Equal(t, false, yaml2Go.Client.Option[5])
	assert.Equal(t, 10, yaml2Go.Client.Option[6])
	assert.Equal(t, 10.011, yaml2Go.Client.Option[7])
	if m, ok := yaml2Go.Client.Option[3].(map[string]interface{}); ok {
		assert.Equal(t, "key_prod", m["sub_options"].(map[string]interface{})["key"])
	} else {
		assert.Fail(t, "expected map")
	}

	// Verify all data - option map
	assert.Equal(t, 3, len(yaml2Go.Client.OptionMap))
	assert.Equal(t, "call_prod", yaml2Go.Client.OptionMap[0]["call"].Name)
	assert.Equal(t, "sms_prod", yaml2Go.Client.OptionMap[1]["sms"].Name)
	assert.Equal(t, "push_prod", yaml2Go.Client.OptionMap[2]["push"].Name)
}

func TestReadParameterizedConfigYamlFile_Prod_List(t *testing.T) {
	fmt.Println()
	_ = os.Setenv("PRDO_testServer", "test.prod")
	_ = os.Setenv("STAGE_testServer", "test.stage")

	type internalClient struct {
		Enabled bool          `yaml:"enabled"`
		Id      string        `yaml:"id"`
		Price   float64       `yaml:"price"`
		Option  []interface{} `yaml:"option"`
	}

	type Yaml2GoInternal struct {
		Client internalClient `yaml:"client"`
	}
	yaml2Go := Yaml2GoInternal{}
	err := ReadParameterizedYaml(strTestReadParameterizedConfigYamlFile_Prod_List, &yaml2Go, "prod")
	assert.NoError(t, err)

	assert.Equal(t, "prod_client", yaml2Go.Client.Id)
	assert.Equal(t, 10.001, yaml2Go.Client.Price)
	assert.Equal(t, true, yaml2Go.Client.Enabled)

	assert.Equal(t, 8, len(yaml2Go.Client.Option))
	assert.Equal(t, "call_prod", yaml2Go.Client.Option[0])
	assert.Equal(t, "sms_prod", yaml2Go.Client.Option[1])
	assert.Equal(t, "email", yaml2Go.Client.Option[2])
	if m, ok := yaml2Go.Client.Option[3].(map[string]interface{}); ok {
		assert.Equal(t, "key_prod", m["sub_options"].(map[string]interface{})["key"])
	} else {
		assert.Fail(t, "expected map")
	}

	assert.Equal(t, true, yaml2Go.Client.Option[4])
	assert.Equal(t, false, yaml2Go.Client.Option[5])
	assert.Equal(t, 10, yaml2Go.Client.Option[6])
	assert.Equal(t, 10.011, yaml2Go.Client.Option[7])

	yaml, err := serialization.ToYaml(yaml2Go)
	assert.NoError(t, err)
	// fmt.Println(yaml)
	_ = yaml
}

func TestReadParameterizedConfigYaml_ListOfMaps(t *testing.T) {
	fmt.Println()
	_ = os.Setenv("PRDO_testServer", "test.prod")
	_ = os.Setenv("STAGE_testServer", "test.stage")

	// ---------------------------------- Wrapper Json Struct ----------------------------------------------------------
	// Internal nodes
	type internalNode struct {
		Name string `yaml:"name"`
	}

	type listMap map[string]internalNode

	// List of options
	type internalList struct {
		Option []listMap `yaml:"option"`
	}

	type wrapper struct {
		Client internalList `yaml:"client"`
	}

	// Read data in struct for env=prod
	yaml2Go := wrapper{}
	err := ReadParameterizedYaml(testStringTestReadParameterizedConfigYamlListOfMaps, &yaml2Go, "prod")
	assert.NoError(t, err)

	// Verify all data
	assert.Equal(t, 3, len(yaml2Go.Client.Option))
	assert.Equal(t, "call_prod", yaml2Go.Client.Option[0]["call"].Name)
	assert.Equal(t, "sms_prod", yaml2Go.Client.Option[1]["sms"].Name)
	assert.Equal(t, "push_prod", yaml2Go.Client.Option[2]["push"].Name)
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
	// fmt.Println(yaml)
	_ = yaml
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
