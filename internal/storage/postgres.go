package storage

import (
	"database/sql"
	"fmt"
	config "system-collector/configs"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient() (*PostgresClient, error) {
	cfg := config.Get()

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres 연결 실패: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres 연결 테스트 실패: %v", err)
	}

	return &PostgresClient{db: db}, nil
}

func (p *PostgresClient) Close() error {
	return p.db.Close()
}

func (p *PostgresClient) GetDB() *sql.DB {
	return p.db
}
