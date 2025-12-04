package main

import (
	"os"
	"strconv"
	"time"
)

func GetenvWithDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var (
	UserAgent = GetenvWithDefault("USER_AGENT", "curl/8.17.0")
	Timeout   = func() time.Duration {
		timeoutEnv := GetenvWithDefault("TIMEOUT", "10")
		timeoutSec, err := strconv.Atoi(timeoutEnv)
		if err != nil {
			panic(err)
		}
		return time.Duration(timeoutSec) * time.Second
	}()
)
