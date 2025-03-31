package repository

import (
	"database/sql"
	"fmt"
	config "system-collector/configs"
	"system-collector/pkg/models"
)

type CommandRepository struct {
	db        *sql.DB
	tableName string
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	cfg := config.Get()
	return &CommandRepository{
		db:        db,
		tableName: cfg.Postgres.TableName,
	}
}

// GetCommandsByNodeID 특정 노드의 명령어 조회
func (r *CommandRepository) GetCommandsByNodeID(nodeID string) ([]models.Command, error) {
	query := fmt.Sprintf(`SELECT command_id, node_id, command_type, command_status 
						 FROM %s WHERE node_id = $1`, r.tableName)
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
	query := fmt.Sprintf(`DELETE FROM %s WHERE node_id = $1`, r.tableName)
	_, err := r.db.Exec(query, nodeID)
	return err
}
