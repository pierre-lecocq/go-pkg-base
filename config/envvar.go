package config

import (
	"fmt"
	"os"
	"strconv"
)

func StringVal(name string) (string, error) {
	strval, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("environment variable %s undefined", name)
	}

	return strval, nil
}

func IntVal(name string) (int, error) {
	strval, found := os.LookupEnv(name)
	if !found {
		return 0, fmt.Errorf("environment variable %s undefined", name)
	}

	intval, err := strconv.Atoi(strval)
	if err != nil {
		return 0, fmt.Errorf("cannot convert environment variable %s to int", name)
	}

	return intval, nil
}
