package serialization

import (
	. "github.com/devlibx/gox-base/v2/errors"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func ReadYaml(file string, object interface{}) (err error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, object)
		if err == nil {
			return nil
		} else {
			return NewError(UnmarshalFailedErrorCode, "could not unmarshal yaml from given file ["+file+"]", err, nil)
		}
	} else {
		return NewError(
			FileOpenErrorCode, "could not open file to read ["+file+"]", err, nil)
	}
}

func ReadYamlFromString(yamlString string, object interface{}) (err error) {
	err = yaml.Unmarshal([]byte(yamlString), object)
	if err == nil {
		return nil
	} else {
		return NewError(UnmarshalFailedErrorCode, "could not unmarshal yaml from given yamlString ["+yamlString+"]", err, nil)
	}
}

func ReadYamlWithEnvVar(file string, object interface{}) (err error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {

		// Resolve env variable
		yamlString := os.ExpandEnv(string(data))
		err = yaml.Unmarshal([]byte(yamlString), object)
		if err == nil {
			return nil
		} else {
			return NewError(UnmarshalFailedErrorCode, "could not unmarshal yaml from given file ["+file+"]", err, nil)
		}
	} else {
		return NewError(
			FileOpenErrorCode, "could not open file to read ["+file+"]", err, nil)
	}
}

func ReadYamlFromStringWithEnvVar(yamlString string, object interface{}) (err error) {
	yamlString = os.ExpandEnv(yamlString)
	return ReadYamlFromString(yamlString, object)
}

func ToYaml(object interface{}) (string, error) {
	data, err := yaml.Marshal(object)
	if err != nil {
		return "", errors.Wrap(err, "failed to write yaml from object")
	}
	return string(data), nil
}

//
