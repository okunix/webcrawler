package config

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

	Timeout = func() time.Duration {
		timeoutEnv := GetenvWithDefault("TIMEOUT", "10")
		timeoutSec, err := strconv.Atoi(timeoutEnv)
		if err != nil {
			panic(err)
		}
		return time.Duration(timeoutSec) * time.Second
	}()

	Crawlers = func() uint32 {
		crawlersEnv := GetenvWithDefault("CRAWLERS", "20")
		crawlers, err := strconv.ParseUint(crawlersEnv, 10, 32)
		if err != nil {
			panic(err)
		}
		return uint32(crawlers)
	}()

	PgHost     = GetenvWithDefault("PG_HOST", "localhost:5432")
	PgUser     = GetenvWithDefault("PG_USER", "postgres")
	PgPassword = GetenvWithDefault("PG_PASSWORD", "postgres")
	PgDb       = GetenvWithDefault("PG_DB", "webcrawler_db")
	PgSSL      = GetenvWithDefault("PG_SSLMODE", "disable")
)
