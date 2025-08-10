package database

import (
	"auction-web/internal/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

// NewPostgresClient connects the application with postgres database and creates new postgres client
func NewPostgresClient(ctx context.Context, uri string) (client *pgx.Conn, err error) {
	client, err = pgx.Connect(ctx, uri)
	if err != nil {
		return nil, logger.WrapError(err, "failed to connect Postgres")
	}

	// Ping the database to verify connection
	if err = client.Ping(ctx); err != nil {
		return nil, logger.WrapError(err, "failed to ping Postgres")
	}

	return client, nil
}

// ExecuteQuery executes a query in postgres
func ExecuteQuery(ctx context.Context, client *pgx.Conn, query string) (err error) {
	if _, err = client.Exec(ctx, query); err != nil {
		return logger.WrapError(err, "failed to execute query")
	}

	return
}

// FetchRecords fetches records from postgres
func FetchRecords(ctx context.Context, client *pgx.Conn, query string, args ...any) (rows pgx.Rows, err error) {
	rows, err = client.Query(ctx, query, args...)
	if err != nil {
		return nil, logger.WrapError(err, "failed to fetch records")
	}

	return rows, nil
}
