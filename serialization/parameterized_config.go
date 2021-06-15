package serialization

import (
	"github.com/devlibx/gox-base/errors"
	"io/ioutil"
	"os"
	"strings"
)

func ReadParameterizedYaml(data string, object interface{}, env string) (err error) {

	// Read all environment var and replace it in input string
	data = os.ExpandEnv(data)

	// Read config as map to resolve parameterized variables
	firstMap := map[string]interface{}{}
	err = ReadYamlFromStringWithEnvVar(data, &firstMap)
	if err != nil {
		return errors.Wrap(err, "could not parse yaml content ["+data+"]", err, nil)
	}

	// Process map in final result map
	newMap := map[string]interface{}{}
	newMap, err = processMap(firstMap, env)
	if err != nil {
		return errors.Wrap(err, "error in parsing yaml content ["+data+"]", err, nil)
	}

	yaml, err := ToYaml(newMap)
	if err != nil {
		return errors.Wrap(err, "could not parse final yaml content ["+data+"]", err, nil)
	}

	return ReadYamlFromString(yaml, object)
}

func ReadParameterizedYamlFile(file string, object interface{}, env string) (err error) {
	_data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "could not open file to read ["+file+"]", err, nil)
	}
	return ReadParameterizedYaml(string(_data), object, env)
}

func processMap(input map[string]interface{}, env string) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for k, v := range input {
		if val, ok := v.(string); ok {
			out[k], _ = processString(val, env)
		} else if val, ok := v.(map[string]interface{}); ok {
			out[k], _ = processMap(val, env)
		} else if val, ok := v.([]interface{}); ok {
			out[k], _ = processList(val, env)
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func processList(input []interface{}, env string) ([]interface{}, error) {
	out := make([]interface{}, 0)
	for _, v := range input {
		if val, ok := v.(string); ok {
			r, _ := processString(val, env)
			out = append(out, r)
		} else if val, ok := v.(map[string]interface{}); ok {
			r, _ := processMap(val, env)
			out = append(out, r)
		} else if val, ok := v.([]interface{}); ok {
			r, _ := processList(val, env)
			out = append(out, r)
		} else {
			out = append(out, v)
		}
	}
	return out, nil
}

func processString(input string, env string) (interface{}, error) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "env:") {
		p := ParameterizedValue(input)
		return p.Get(env)
	} else {
		return input, nil
	}
}
