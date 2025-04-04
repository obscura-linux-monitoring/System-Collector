package repository

import (
	"database/sql"
	"system-collector/pkg/logger"
	"system-collector/pkg/models"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	sugar := logger.GetCustomLogger()
	sugar.Infow("CommandRepository 초기화 중")

	return &CommandRepository{
		db: db,
	}
}

// GetCommandsByNodeID 특정 노드의 명령어 조회
func (r *CommandRepository) GetCommandsByNodeID(nodeID string) ([]models.Command, error) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("노드의 명령어 조회 시작", "nodeID", nodeID)

	query := `SELECT command_id, node_id, command_type, command_status, target FROM commands WHERE node_id = $1`
	rows, err := r.db.Query(query, nodeID)
	if err != nil {
		sugar.Errorw("명령어 조회 SQL 오류", "nodeID", nodeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var commands []models.Command
	for rows.Next() {
		var cmd models.Command
		if err := rows.Scan(
			&cmd.CommandID,
			&cmd.NodeID,
			&cmd.CommandType,
			&cmd.CommandStatus,
			&cmd.Target,
		); err != nil {
			sugar.Errorw("명령어 데이터 스캔 오류", "error", err)
			return nil, err
		}
		commands = append(commands, cmd)
	}

	sugar.Infow("명령어 조회 완료", "nodeID", nodeID, "commandCount", len(commands))
	return commands, nil
}

// deleteCommandsByNodeID 특정 노드의 모든 명령어 삭제
func (r *CommandRepository) DeleteCommandsByNodeID(nodeID string) error {
	sugar := logger.GetCustomLogger()
	sugar.Infow("노드의 명령어 삭제 시작", "nodeID", nodeID)

	query := `DELETE FROM commands WHERE node_id = $1`
	_, err := r.db.Exec(query, nodeID)
	if err != nil {
		sugar.Errorw("명령어 삭제 SQL 오류", "nodeID", nodeID, "error", err)
		return err
	}

	sugar.Infow("노드의 명령어 삭제 완료", "nodeID", nodeID)
	return nil
}
