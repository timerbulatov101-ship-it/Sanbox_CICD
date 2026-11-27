package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Store struct {
	db *sql.DB
}

func New(ctx context.Context, cfg Config) (*Store, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("incomplete PostgreSQL configuration")
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}
	query := dsn.Query()
	query.Set("sslmode", cfg.SSLMode)
	dsn.RawQuery = query.Encode()

	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS orders (
		id BIGSERIAL PRIMARY KEY,
		order_id BIGINT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("cannot create orders table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) SaveOrder(ctx context.Context, orderID int64) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO orders (order_id) VALUES ($1)", orderID)
	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}
