package repository

import (
	"database/sql"
	"system-collector/pkg/models"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	return &CommandRepository{
		db: db,
	}
}

// GetCommandsByNodeID 특정 노드의 명령어 조회
func (r *CommandRepository) GetCommandsByNodeID(nodeID string) ([]models.Command, error) {
	query := `SELECT command_id, node_id, command_type, command_status FROM commands WHERE node_id = $1`
	rows, err := r.db.Query(query, nodeID)
	if err != nil {
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
		); err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

// deleteCommandsByNodeID 특정 노드의 모든 명령어 삭제
func (r *CommandRepository) DeleteCommandsByNodeID(nodeID string) error {
	query := `DELETE FROM commands WHERE node_id = $1`
	_, err := r.db.Exec(query, nodeID)
	return err
}
