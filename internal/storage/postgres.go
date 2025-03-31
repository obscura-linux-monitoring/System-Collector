package storage

import (
	"database/sql"
	"fmt"
	config "system-collector/configs"
	"system-collector/pkg/logger"

	_ "github.com/lib/pq"
)

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient() (*PostgresClient, error) {
	sugar := logger.GetSugar()
	sugar.Info("PostgresClient 초기화 중")

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
		sugar.Errorw("postgres 연결 실패", "error", err)
		return nil, fmt.Errorf("postgres 연결 실패: %v", err)
	}

	if err := db.Ping(); err != nil {
		sugar.Errorw("postgres 연결 테스트 실패", "error", err)
		return nil, fmt.Errorf("postgres 연결 테스트 실패: %v", err)
	}

	sugar.Info("postgres 연결 성공")
	return &PostgresClient{db: db}, nil
}

func (p *PostgresClient) Close() error {
	sugar := logger.GetSugar()
	sugar.Info("postgres 연결 종료 중")

	err := p.db.Close()
	if err != nil {
		sugar.Errorw("postgres 연결 종료 실패", "error", err)
		return fmt.Errorf("postgres 연결 종료 실패: %v", err)
	}

	sugar.Info("postgres 연결 종료 완료")
	return nil
}

func (p *PostgresClient) GetDB() *sql.DB {
	sugar := logger.GetSugar()
	sugar.Info("postgres 데이터베이스 연결 반환")

	return p.db
}
