package repository

import (
	"database/sql"
	"system-collector/pkg/logger"
	"system-collector/pkg/models"
)

type LogRepository struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	sugar := logger.GetCustomLogger()
	sugar.Infof("LogRepository 초기화 중")

	return &LogRepository{
		db: db,
	}
}

func (r *LogRepository) SaveLogs(logs []models.LogMessage) error {
	sugar := logger.GetCustomLogger()
	sugar.Infof("로그 저장 시작 %d개", len(logs))

	query := `INSERT INTO logs (node_id, timestamp, level, content) VALUES ($1, $2, $3, $4)`
	for _, log := range logs {
		_, err := r.db.Exec(query, log.NodeID, log.Timestamp, log.Level, log.Content)
		if err != nil {
			sugar.Errorw("로그 저장 실패", "error", err)
			return err
		}
	}

	sugar.Infof("로그 저장 완료 %d개", len(logs))
	return nil
}
