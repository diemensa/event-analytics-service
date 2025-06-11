package config

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitPostgres(host, user, pw, dbname, port string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pw, dbname, port,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	createTableQuery :=
		`CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY,
		type TEXT NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL,
		user_id UUID
	);`

	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return db, nil
}
