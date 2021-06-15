package config

import (
	"github.com/devlibx/gox-base/errors"
	"strconv"
	"strings"
)

type ParameterizedValue string

type parseFunction func(in string) (interface{}, string, error)

func (p ParameterizedValue) GetInt(env string) (int, error) {
	if val, err := p.Get(env); err != nil {
		return 0, err
	} else if finalValue, ok := val.(int); ok {
		return finalValue, nil
	} else if str, ok := val.(string); ok {
		return strconv.Atoi(str)
	} else {
		return 0, errors.New("not a int: value=%v", val)
	}
}

func (p ParameterizedValue) GetString(env string) (string, error) {
	if val, err := p.Get(env); err != nil {
		return "", err
	} else if finalValue, ok := val.(string); ok {
		return finalValue, nil
	} else {
		return "", errors.New("not a string: value=%v", val)
	}
}

func (p ParameterizedValue) GetBool(env string) (bool, error) {
	if val, err := p.Get(env); err != nil {
		return false, err
	} else if finalValue, ok := val.(bool); ok {
		return finalValue, nil
	} else if str, ok := val.(string); ok {
		return strconv.ParseBool(str)
	} else {
		return false, errors.New("not a bool: value=%v", val)
	}
}

func (p ParameterizedValue) GetFloat(env string) (float64, error) {
	if val, err := p.Get(env); err != nil {
		return 0, err
	} else if finalValue, ok := val.(float64); ok {
		return finalValue, nil
	} else if str, ok := val.(string); ok {
		return strconv.ParseFloat(str, 64)
	} else {
		return 0, errors.New("not a float: value=%v", val)
	}
}

func (p ParameterizedValue) Get(env string) (interface{}, error) {
	str := strings.TrimSpace(string(p))
	var pf parseFunction = func(in string) (interface{}, string, error) {
		return in, "string", nil
	}

	// If this is not parametrized then just return the value
	if !strings.HasPrefix(str, "env:") {
		return str, nil
	}

	if strings.HasPrefix(str, "env:string:") {
		str = strings.Replace(str, "env:string:", "", 1)
		pf = func(in string) (interface{}, string, error) {
			return in, "string", nil
		}
	} else if strings.HasPrefix(str, "env:int:") {
		str = strings.Replace(str, "env:int:", "", 1)
		pf = func(in string) (interface{}, string, error) {
			val, err := strconv.Atoi(in)
			return val, "int", err
		}
	} else if strings.HasPrefix(str, "env:bool:") {
		str = strings.Replace(str, "env:bool:", "", 1)
		pf = func(in string) (interface{}, string, error) {
			val, err := strconv.ParseBool(in)
			return val, "bool", err
		}
	} else if strings.HasPrefix(str, "env:float:") {
		str = strings.Replace(str, "env:float:", "", 1)
		pf = func(in string) (interface{}, string, error) {
			val, err := strconv.ParseFloat(in, 64)
			return val, "float", err
		}
	}

	tokens := strings.Split(str, ";")

	// Find matching env string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return false, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == env {
			part := strings.TrimSpace(parts[1])
			if val, tname, err := pf(part); err == nil {
				return val, err
			} else {
				return false, errors.Wrap(err, "expected the value to be a %s buf got [%s]", tname, str)
			}
		}
	}

	// Find default string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return false, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == "default" {
			part := strings.TrimSpace(parts[1])
			if val, tname, err := pf(part); err == nil {
				return val, err
			} else {
				return false, errors.Wrap(err, "expected the value to be a %s buf got [%s]", tname, str)
			}
		}
	}

	return nil, errors.New("did not find value: str=[%s] env=%s", p, env)
}
