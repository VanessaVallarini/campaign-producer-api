package postgres

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-producer-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	pgxDD "gopkg.in/DataDog/dd-trace-go.v1/contrib/jackc/pgx.v5"
)

// buildConnString constructs the database connection string from the config.
func buildConnString(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s default_query_exec_mode=cache_describe",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Database,
		cfg.Password,
		"disable",
		cfg.Conn.Min,
		cfg.Conn.Max,
		cfg.Conn.Lifetime,
		cfg.Conn.IdleTime,
	)
}

// CreatePool initializes and returns a new pgxpool.Pool using the given configuration.
func CreatePool(ctx context.Context, cfg *config.DatabaseConfig) *pgxpool.Pool {
	connString := buildConnString(cfg)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		err := errors.Wrapf(err, "Unable to parse config: %s", connString)
		easyzap.Fatal(ctx, err, "unable to parse config")

		return nil
	}

	pool, err := pgxDD.NewPoolWithConfig(ctx, poolConfig)
	if err != nil {
		err := errors.Wrapf(err, "Unable to create pool with config: %s", connString)
		easyzap.Fatal(ctx, err, "nable to create pool with config")

		return nil
	}

	easyzap.Info(ctx, "database connection pool created successfully.")

	return pool
}
