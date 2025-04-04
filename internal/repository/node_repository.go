package repository

import (
	"database/sql"
	"system-collector/pkg/logger"
	"system-collector/pkg/models"
)

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	sugar := logger.GetCustomLogger()
	sugar.Infow("NodeRepository 초기화 중")

	return &NodeRepository{
		db: db,
	}
}

func (r *NodeRepository) CreateNode(node *models.Node) error {
	sugar := logger.GetCustomLogger()
	sugar.Infow("노드 생성 시작", "node", node)

	query := `INSERT INTO nodes (node_id, obscura_key, server_type) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, node.NodeID, node.ObscuraKey, node.ServerType)
	if err != nil {
		sugar.Errorw("노드 생성 실패", "error", err)
		return err
	}
	return err
}

func (r *NodeRepository) GetAllNodes() ([]*models.Node, error) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("모든 노드 조회 시작")

	query := `SELECT node_id, obscura_key, server_type FROM nodes`
	rows, err := r.db.Query(query)
	if err != nil {
		sugar.Errorw("모든 노드 조회 실패", "error", err)
		return nil, err
	}
	defer rows.Close()

	nodes := []*models.Node{}
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.NodeID, &node.ObscuraKey, &node.ServerType)
		if err != nil {
			sugar.Errorw("노드 데이터 스캔 오류", "error", err)
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	sugar.Infow("모든 노드 조회 완료", "nodeCount", len(nodes))
	return nodes, nil
}
