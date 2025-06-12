package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPostgres(ctx context.Context, host, user, pw, dbname, port string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pw, dbname, port,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	createTableQuery :=
		`CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY,
		type TEXT NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL,
		user_id UUID
	);`

	if _, err = pool.Exec(ctx, createTableQuery); err != nil {
		return nil, err
	}

	return pool, nil
}
