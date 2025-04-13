package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// initDB connects to Postgres and returns a *sqlx.DB.
func initDB() (*sqlx.DB, error) {
	// Update these creds as needed:
	user, pass, host, port, name := "user", "pass", "localhost", 5432, "productsdb"
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, name,
	)
	return sqlx.Connect("postgres", dsn)
}
