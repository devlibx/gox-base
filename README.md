Gox-Base project provide utilities which is used commonly in all applications.

1. Serialization utils
2. Json file to object
3. Yaml file to object
4. XML file to object
5. ...

#Config
You can use this for common application configuration
```go
type configToTest struct {
	App    App    `json:"app" yaml:"app"`
	Logger Logger `json:"logger" yaml:"logger"`
}

// Read struct with YAML
conf := &configToTest{}
err := serialization.ReadYamlFromString(yamlConfig, conf)
assert.NoError(t, err)


// Read struct with Json
conf := &configToTest{}
err := serialization.JsonToObject(jsonString, conf)
assert.NoError(t, err)
```

You can also extend an existing type to include more. For example, we will add one more property called
```some``` in ```App``` type
```go

// Extended "App" object with one more type "Some"
// NOTE - in case you have YAM as source you must put `yaml:",inline"`
type extendedAppObject struct {
	 App 		 `yaml:",inline"`
	 Some string `json:"some"`
}

type configToTestWithExtendedAppObject struct {
	App    extendedAppObject `json:"app" yaml:"app"`
	Logger Logger            `json:"logger" yaml:"logger"`
}
```



# Utility

##### Setup cross function which is usd in almost all apis to be used in gox

```go
import (
    "github.com/devlibx/gox-base"
    "go.uber.org/zap"
)

zapConfig := zap.NewDevelopmentConfig()
zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
crossFunction := gox.NewCrossFunction(zapConfig.Build())

For test cases:
===============
import "github.com/devlibx/gox-base/test"
cf, gomockController := test.MockCf(t, zap.InfoLevel)
```

#### Convert anything to string

This utility will convert int, bool, interface{} to string. Object will output a json string.

```golang
out, _ := Stringify(10)
// Output = "10"

boolOut, _ := Stringify(true)
// Output = "true"

type utilTestStruct struct {
IntValue    int    `json:"int"`
BoolValue   bool   `json:"bool"`
StringValue string `json:"string"`
}

objectOut, _ := Stringify(utilTestStruct{
IntValue:    10,
BoolValue:   false,
StringValue: "some value",
})
// Output = {"int":10,"bool":false,"string":"some value"}
```

###### Stringify with error suppressed

If you don't want to handle error and have default value on error then you can use "suppress error" version of this
method

```golang
intOut1 := StringifySuppressError(10, "0")
// Output = "10"

If there is a error when input is bad then you will get the default
value "0"
```

### Logging

You can use logger from Uber

```go
var cfg zap.Config
cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
// cfg.Encoding = "console"
cfg.Encoding = "json"
cfg.OutputPaths = []string{"stdout", "/tmp/logs1"}
cfg.ErrorOutputPaths = []string{"stdout", "/tmp/logs2"}
cfg.EncoderConfig = zapcore.EncoderConfig{
MessageKey:     "message",
LevelKey:       "level",
EncodeLevel:    zapcore.LowercaseLevelEncoder,
EncodeDuration: zapcore.SecondsDurationEncoder,
EncodeCaller:   zapcore.ShortCallerEncoder,
StacktraceKey:  "stacktrace",
TimeKey:        "timestamp",
EncodeTime:     zapcore.ISO8601TimeEncoder,
}

// Production encoder
// encoderCfg := zap.NewProductionEncoderConfig()
// encoderCfg.TimeKey = "timestamp"
// encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
// cfg.EncoderConfig = encoderCfg

logger, _ := cfg.Build()
defer logger.Sync()

logger.Debug("Logger from parent level - it has not key-value")

// A module level logger (you can log user Id or any param here)
subModuleLogger := logger.With(zap.String("userId", "1234"))
subModuleLogger.Debug("this is a logger for sub-module")

logger.Debug("Logger from parent level - it has not key-value")
```

Some more example https://github.com/uber-go/zap/blob/master/example_test.go
