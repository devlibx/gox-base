package serialization

import (
	. "github.com/harishb2k/gox-base/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ReadYaml(file string, object interface{}) (err error) {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(data, object)
		if err == nil {
			return nil
		} else {
			return NewError(UnmarshalFailedErrorCode, "could not unmarshal json from given file ["+file+"]", err, nil)
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
		return NewError(UnmarshalFailedErrorCode, "could not unmarshal json from given yamlString ["+yamlString+"]", err, nil)
	}
}
