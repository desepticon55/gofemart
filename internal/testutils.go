package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"path/filepath"
	"testing"
	"time"
)

func InitPostgresIntegrationTest(t *testing.T, ctx context.Context, logger *zap.Logger) (*pgxpool.Pool, func() error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	runMigrations(connStr, logger)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	assert.NoError(t, err)
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	cleanup := func() error {
		if err := pgContainer.Terminate(ctx); err != nil {
			return fmt.Errorf("failed to terminate pgContainer: %w", err)
		}
		return nil
	}

	return pool, cleanup
}

func ClearTables(ctx context.Context, pool *pgxpool.Pool) error {
	tables := []string{"balance", "withdrawal", "order", "user"}
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE gofemart.%s CASCADE", table)
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}

func runMigrations(connectionString string, logger *zap.Logger) {
	databaseConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		logger.Error("Error during parse database URL", zap.Error(err))
		return
	}
	db := stdlib.OpenDB(*databaseConfig)
	defer db.Close()

	goose.SetDialect("postgres")
	if err := goose.Up(db, filepath.Join("../..", "migrations")); err != nil {
		logger.Error("Error during run database migrations", zap.Error(err))
	}
}
