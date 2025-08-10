package database

import (
	"auction-web/internal/constants"
	"auction-web/internal/logger"
	"context"
	"math"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgresClient connects the application with postgres database and creates new postgres client
func NewPostgresClient(ctx context.Context, uri string) (client *pgxpool.Pool, err error) {
	retries := 1
	for {
		client, err = pgxpool.Connect(ctx, uri)
		if err != nil {
			if retries > constants.MaxRetries {
				return nil, logger.WrapError(err, "failed to connect Postgres")
			}
			time.Sleep(time.Duration(math.Pow(2, float64(retries))) * time.Second)
			continue
		}
		break
	}

	// Ping the database to verify connection
	if err = pingDB(ctx, client); err != nil {
		return nil, logger.WrapError(err, "failed to ping Postgres")
	}

	return client, nil
}

// pingDB pings to the postgres database
func pingDB(ctx context.Context, client *pgxpool.Pool) error {
	conn, err := client.Acquire(ctx)
	if err != nil {
		return logger.WrapError(err, "failed to acquire connection")
	}

	// Ping the database
	if err = conn.Conn().Ping(ctx); err != nil {
		return logger.WrapError(err, "failed to ping database")
	}

	return nil
}

// ExecuteQuery executes a query in postgres
func ExecuteQuery(ctx context.Context, client *pgxpool.Pool, query string, args ...any) (err error) {
	if _, err = client.Exec(ctx, query, args...); err != nil {
		return logger.WrapError(err, "failed to execute query")
	}

	return
}

func FetchRecords[T any](ctx context.Context, client *pgxpool.Pool, query string, args ...any) (results []T, err error) {
	if err = pgxscan.Select(ctx, client, &results, query, args...); err != nil {
		return nil, logger.WrapError(err, "failed to fetch records")
	}

	return results, nil
}
