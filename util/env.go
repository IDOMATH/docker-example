package util

import (
	"errors"
	"os"
)

func EnvMust(envValue string) (string, error) {
	str := os.Getenv(envValue)
	if str == "" {
		return "", errors.New("env variable not found")
	}
	return str, nil
}
