package config

import (
	"github.com/devlibx/gox-base/errors"
	"strconv"
	"strings"
)

type ParameterizedValue string

type ParameterizedString string
type ParameterizedInt string
type ParameterizedBool string
type ParameterizedFloat string

func (p ParameterizedString) Get(env string) (string, error) {

	// See if this is a parameterized string or not - if not then just give back the data
	str := string(p)
	if !strings.HasPrefix(str, "env:") {
		return str, nil
	}

	str = strings.Replace(str, "env:", "", 1)
	tokens := strings.Split(str, ";")

	// Find matching env string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return "", errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == env {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	// Find default string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return "", nil
		} else if strings.TrimSpace(parts[0]) == "default" {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", errors.New("env not found")
}

func (p ParameterizedInt) Get(env string) (int, error) {

	// If it is already a integer then just return it
	str := string(p)
	if val, err := strconv.Atoi(str); err == nil {
		return val, err
	}

	// See if this is a parameterized string or not - if not then just give back the data
	if !strings.HasPrefix(str, "env:") {
		if val, err := strconv.Atoi(str); err == nil {
			return val, err
		} else {
			return 0, err
		}
	}

	str = strings.Replace(str, "env:", "", 1)
	tokens := strings.Split(str, ";")

	// Find matching env string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return 0, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == env {
			if val, err := strconv.Atoi(parts[1]); err == nil {
				return val, err
			} else {
				return 0, errors.Wrap(err, "expected the value to be a int buf got [%s]", str)
			}
		}
	}

	// Find default string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return 0, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == "default" {
			if val, err := strconv.Atoi(parts[1]); err == nil {
				return val, err
			} else {
				return 0, errors.Wrap(err, "expected the value to be a int buf got [%s]", str)
			}
		}
	}

	return 0, errors.New("env not found")
}

func (p ParameterizedFloat) Get(env string) (float64, error) {

	// If it is already a integer then just return it
	str := string(p)
	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val, err
	}

	// See if this is a parameterized string or not - if not then just give back the data
	if !strings.HasPrefix(str, "env:") {
		if val, err := strconv.ParseFloat(str, 64); err == nil {
			return val, err
		} else {
			return 0, err
		}
	}

	str = strings.Replace(str, "env:", "", 1)
	tokens := strings.Split(str, ";")

	// Find matching env string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return 0, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == env {
			if val, err := strconv.ParseFloat(parts[1], 64); err == nil {
				return val, err
			} else {
				return 0, errors.Wrap(err, "expected the value to be a int buf got [%s]", str)
			}
		}
	}

	// Find default string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return 0, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == "default" {
			if val, err := strconv.ParseFloat(parts[1], 64); err == nil {
				return val, err
			} else {
				return 0, errors.Wrap(err, "expected the value to be a int buf got [%s]", str)
			}
		}
	}

	return 0, errors.New("env not found")
}

func (p ParameterizedBool) Get(env string) (bool, error) {

	// If it is already a integer then just return it
	str := string(p)
	if val, err := strconv.ParseBool(str); err == nil {
		return val, err
	}

	// See if this is a parameterized string or not - if not then just give back the data
	if !strings.HasPrefix(str, "env:") {
		if val, err := strconv.ParseBool(str); err == nil {
			return val, err
		} else {
			return false, err
		}
	}

	str = strings.Replace(str, "env:", "", 1)
	tokens := strings.Split(str, ";")

	// Find matching env string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return false, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == env {
			if val, err := strconv.ParseBool(parts[1]); err == nil {
				return val, err
			} else {
				return false, errors.Wrap(err, "expected the value to be a bool buf got [%s]", str)
			}
		}
	}

	// Find default string
	for _, token := range tokens {
		parts := strings.Split(token, "=")
		if len(parts) < 2 {
			return false, errors.New("expected a key=value but got [%s]", token)
		} else if strings.TrimSpace(parts[0]) == "default" {
			if val, err := strconv.ParseBool(parts[1]); err == nil {
				return val, err
			} else {
				return false, errors.Wrap(err, "expected the value to be a bool buf got [%s]", str)
			}
		}
	}

	return false, errors.New("env not found")
}

type parseFunction func(in string) (interface{}, string, error)

func (p ParameterizedValue) Get(env string) (interface{}, error) {
	str := strings.TrimSpace(string(p))
	var pf parseFunction = func(in string) (interface{}, string, error) {
		return in, "string", nil
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

	return nil, nil
}
