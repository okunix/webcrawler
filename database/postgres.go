package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/okunix/webcrawler/config"
)

var (
	initPostgres sync.Once
	postgres     *sql.DB
)

func Postgres() *sql.DB {
	initPostgres.Do(func() {
		addr := fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=%s",
			config.PgUser,
			config.PgPassword,
			config.PgHost,
			config.PgDb,
			config.PgSSL,
		)
		conn, err := sql.Open("postgres", addr)
		if err != nil {
			panic(err)
		}
		if err := conn.Ping(); err != nil {
			panic(err)
		}
		postgres = conn
	})
	return postgres
}
