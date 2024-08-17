package serialization

import (
	"encoding/json"
	"fmt"
	"github.com/devlibx/gox-base/v2/util"
	"strconv"
)

// Convert input to a string (object is converted to Json String)
// This method will suppress error and will give default value on error
func StringifyOrEmptyJsonOnError(input interface{}) (out string) {
	if input == nil {
		return "{}"
	}
	if out, err := Stringify(input); err == nil {
		if util.IsStringEmpty(out) {
			return "{}"
		} else {
			return out
		}
	} else {
		return "{}"
	}
}

// Convert input to a string (object is converted to Json String)
// This method will suppress error and will give default value on error
func StringifyOrDefaultOnError(input interface{}, valueOnError string) (out string) {
	if out, err := Stringify(input); err == nil {
		return out
	} else {
		return valueOnError
	}
}

// Convert input to a string (object is converted to Json String)
// This method will suppress error and will give default value on error
func StringifyOrEmptyOnError(input interface{}) (out string) {
	if input == nil {
		return ""
	}
	if out, err := Stringify(input); err == nil {
		if util.IsStringEmpty(out) {
			return ""
		} else {
			return out
		}
	} else {
		return ""
	}
}

// Convert input to a string (object is converted to Json String)
// This method will suppress error and will give default value on error
func StringifySuppressError(input interface{}, valueOnError string) (out string) {
	if out, err := Stringify(input); err == nil {
		return out
	} else {
		return valueOnError
	}
}

// Convert input to a string (object is converted to Json String)
func Stringify(input interface{}) (out string, err error) {
	switch v := input.(type) {

	case int:
		out = strconv.Itoa(v)

	case int8, int16, int32, int64:
		out = fmt.Sprintf("%d", v)

	case bool:
		if v {
			out = "true"
		} else {
			out = "false"
		}

	case string:
		out = v

	case []byte:
		out = string(v)

	default:
		if _out, err := json.Marshal(v); err != nil {
			return "", err
		} else {
			out = string(_out)
		}
	}
	return
}
