package database

import (
	"auction-web/internal/constants"
	"auction-web/internal/logger"
	"context"
	"math"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

// NewPostgresClient connects the application with postgres database and creates new postgres client
func NewPostgresClient(ctx context.Context, uri string) (client *pgxpool.Pool, err error) {
	var retries = 1

	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, logger.WrapError(err, "invalid postgres uri")
	}

	// Minimal memory, fast response tuning
	cfg.MinConns = 1                      // keep 1 conn always open
	cfg.MaxConns = 4                      // atmost 4 query can execute parallely
	cfg.MaxConnIdleTime = 5 * time.Minute // close unused conn after 5m
	cfg.MaxConnLifetime = time.Hour       // recycle every hour
	cfg.HealthCheckPeriod = time.Minute   // check idle conns every 1m
	cfg.LazyConnect = false               // connect immediately at startup

	for {
		client, err = pgxpool.ConnectConfig(ctx, cfg)
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
		client.Close()
		return nil, logger.WrapError(err, "failed to ping Postgres")
	}

	// Start logging stats in background 
	// TODO Will uncomment in production
	// go logPoolStats(ctx, client)

	return client, nil
}

// pingDB pings to the postgres database
func pingDB(ctx context.Context, client *pgxpool.Pool) error {
	conn, err := client.Acquire(ctx)
	if err != nil {
		return logger.WrapError(err, "failed to acquire connection")
	}
	defer conn.Release()

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

	return nil
}

// ExecuteQueryReturning executes a query in postgres and return the needed value
func ExecuteQueryReturning(ctx context.Context, client *pgxpool.Pool, dest any, query string, args ...any) (err error) {
	if err = client.QueryRow(ctx, query, args...).Scan(dest); err != nil {
		return logger.WrapError(err, "failed to execute returning query")
	}

	return nil
}

// FetchRecords fetches the rows from postgres
func FetchRecords[T any](ctx context.Context, client *pgxpool.Pool, query string, args ...any) (results []T, err error) {
	if err = pgxscan.Select(ctx, client, &results, query, args...); err != nil {
		return nil, logger.WrapError(err, "failed to fetch records")
	}

	return results, nil
}

// logPoolStats logs connection pool statistics periodically.
func logPoolStats(ctx context.Context, pool *pgxpool.Pool) {
	log := logger.Get()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats := pool.Stat()
			log.Info("Postgres pool stats",
				zap.Int("acquired", int(stats.AcquiredConns())), // in use
				zap.Int("idle", int(stats.IdleConns())),         // ready but unused
				zap.Int("total", int(stats.TotalConns())),       // all connections
				zap.Int32("maxConns", pool.Config().MaxConns),
			)
		}
	}
}
