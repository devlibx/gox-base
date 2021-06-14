package config

import (
	"github.com/devlibx/gox-base/errors"
	"github.com/devlibx/gox-base/serialization"
	"io/ioutil"
	"os"
	"strings"
)

func ReadParameterizedYaml(data string, object interface{}, env string) (err error) {

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) >= 2 {
			data = strings.ReplaceAll(data, "$"+pair[0], pair[1])
		}
	}

	firstMap := map[string]interface{}{}
	err = serialization.ReadYamlFromStringWithEnvVar(data, &firstMap)
	if err != nil {
		return errors.Wrap(err, "could not parse yaml content ["+data+"]", err, nil)
	}

	newMap := map[string]interface{}{}
	processMap(firstMap, newMap, env)

	yaml, err := serialization.ToYaml(newMap)
	if err != nil {
		return errors.Wrap(err, "could not parse final yaml content ["+data+"]", err, nil)
	} else {
		// fmt.Println(yaml)
	}

	return serialization.ReadYamlFromString(yaml, object)
}

func ReadParameterizedYamlFile(file string, object interface{}, env string) (err error) {
	_data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "could not open file to read ["+file+"]", err, nil)
	}
	return ReadParameterizedYaml(string(_data), object, env)
}

func processMap(m map[string]interface{}, newMap map[string]interface{}, env string) {
	for key, value := range m {
		if val, ok := value.(map[string]interface{}); ok {
			subMap := map[string]interface{}{}
			processMap(val, subMap, env)
			newMap[key] = subMap
		} else if val, ok := value.(string); ok {
			if strings.HasPrefix(val, "env:string:") {
				val = strings.Replace(val, "env:string:", "env:", 1)
				p := ParameterizedString(val)
				newMap[key], _ = p.Get(env)
			} else if strings.HasPrefix(val, "env:bool:") {
				val = strings.Replace(val, "env:bool:", "env:", 1)
				p := ParameterizedBool(val)
				newMap[key], _ = p.Get(env)
			} else if strings.HasPrefix(val, "env:int:") {
				val = strings.Replace(val, "env:int:", "env:", 1)
				p := ParameterizedInt(val)
				newMap[key], _ = p.Get(env)
			} else if strings.HasPrefix(val, "env:float:") {
				val = strings.Replace(val, "env:float:", "env:", 1)
				p := ParameterizedFloat(val)
				newMap[key], _ = p.Get(env)
			} else {
				newMap[key] = val
			}
		} else if val, ok := value.([]interface{}); ok && false { // FIXME - we need to handle this
			subList := make([]interface{}, 0)
			for _, v := range val {
				_ = v
			}
			newMap[key] = subList
		} else {
			newMap[key] = value
		}
	}
}
