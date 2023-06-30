package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
)

func NewConnectionDB(ctx context.Context, c *config.DbConfig) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s", c.Host, c.Port, c.DBName, c.SSLMode)
	conn, err := pgx.Connect(ctx, connStr)

	if c == nil {
		logrus.Fatalf("failed to connect the database: %v", err)
		return nil, err
	}
	if err != nil {
		logrus.Fatalf("failed to connect to the database: %v", err)
		return nil, err
	}

	if err = conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	logrus.Println("database connection successfully")

	return conn, nil
}
