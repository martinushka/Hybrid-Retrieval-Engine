package product

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectPostgres(ctx context.Context, databaseURL string) (*pgx.Conn, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return conn, nil
}
