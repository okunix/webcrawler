package main

import "os"

func GetenvWithDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var (
	UserAgent = GetenvWithDefault("USER_AGENT", "curl/8.17.0")
	Timeout   = GetenvWithDefault("TIMEOUT", "10")
)
