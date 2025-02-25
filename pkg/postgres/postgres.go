package postgres

import (
	"Storage/internal/config"
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func NewPostgres(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		strconv.Itoa(cfg.Postgres.Port),
		cfg.Postgres.DB)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("New Postgre: %w", err)
	}
	return conn, nil
}
