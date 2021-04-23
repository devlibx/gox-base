package serialization

import (
	"bytes"
	"encoding/json"
	. "github.com/divlibx/gox-base/errors"
	"io/ioutil"
)

// Read a json file and populate the given object with its content
func ReadJson(file string, object interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err == nil {
		err = json.Unmarshal(data, object)
		if err != nil {
			return NewError(UnmarshalFailedErrorCode, "could not unmarshal json from given file ["+file+"]", err, nil)
		}
	} else {
		return NewError(FileOpenErrorCode, "could not open file to read ["+file+"]", err, nil)
	}
	return nil
}

func ToBytes(object interface{}) ([]byte, error) {

	// Nothing is done if input is nil
	if object == nil {
		return nil, nil
	}

	// If it is already a byte array then return it - no work required
	if obj, ok := object.([]byte); ok {
		return obj, nil
	}

	// Try to read object and convert to []byte representing a json
	response := new(bytes.Buffer)
	if err := json.NewEncoder(response).Encode(object); err != nil {
		return nil, Wrap(err, "fail to write json object to byte array")
	} else {
		return response.Bytes(), nil
	}
}

func ToBytesSuppressError(object interface{}) []byte {
	data, _ := ToBytes(object)
	return data
}

// Read string and fill up the input object
func JsonToObject(input string, object interface{}) error {
	return json.Unmarshal([]byte(input), object)
}

// Helper to convert byte data to a object
func JsonBytesToObject(input []byte, out interface{}) (err error) {
	return json.Unmarshal(input, out)
}

// Helper to convert byte data to a object
func JsonBytesToObjectSuppressError(input []byte, out interface{}) {
	_ = json.Unmarshal(input, out)
}
