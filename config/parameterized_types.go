package config

import (
	"github.com/devlibx/gox-base/errors"
	"strconv"
	"strings"
)

type ParameterizedString string
type ParameterizedInt string
type ParameterizedBool string

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
