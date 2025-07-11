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

	query := `INSERT INTO nodes (node_id, obscura_key, server_type, node_name) VALUES ($1, $2, $3, 'Default')`
	_, err := r.db.Exec(query, node.NodeID, node.ObscuraKey, node.ServerType)
	if err != nil {
		sugar.Errorw("노드 생성 실패", "error", err)
		return err
	}
	return nil
}

func (r *NodeRepository) GetAllNodes() ([]*models.Node, error) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("모든 노드 조회 시작")

	query := `SELECT node_id, obscura_key, server_type, COALESCE(external_ip, '') FROM nodes`
	rows, err := r.db.Query(query)
	if err != nil {
		sugar.Errorw("모든 노드 조회 실패", "error", err)
		return nil, err
	}
	defer rows.Close()

	nodes := []*models.Node{}
	for rows.Next() {
		var node models.Node
		err := rows.Scan(&node.NodeID, &node.ObscuraKey, &node.ServerType, &node.ExternalIP)
		if err != nil {
			sugar.Errorw("노드 데이터 스캔 오류", "error", err)
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	sugar.Infow("모든 노드 조회 완료", "nodeCount", len(nodes))
	return nodes, nil
}

func (r *NodeRepository) UpdateNodeStatus(nodeID string, status int) error {
	sugar := logger.GetCustomLogger()
	sugar.Infof("노드 상태 업데이트: %s, %d", nodeID, status)

	query := `UPDATE nodes SET status = $1 WHERE node_id = $2`
	_, err := r.db.Exec(query, status, nodeID)
	if err != nil {
		sugar.Errorw("노드 상태 업데이트 실패", "error", err)
		return err
	}
	return nil
}

// UpdateNodeExternalIP는 노드의 외부 IP 주소를 업데이트합니다.
func (r *NodeRepository) UpdateNodeExternalIP(nodeID, externalIP string) error {
	sugar := logger.GetCustomLogger()
	sugar.Infow("노드 외부 IP 업데이트", "nodeID", nodeID, "externalIP", externalIP)

	query := `UPDATE nodes SET external_ip = $1 WHERE node_id = $2`
	_, err := r.db.Exec(query, externalIP, nodeID)
	if err != nil {
		sugar.Errorw("노드 외부 IP 업데이트 실패", "error", err, "nodeID", nodeID, "externalIP", externalIP)
		return err
	}

	sugar.Infow("노드 외부 IP 업데이트 성공", "nodeID", nodeID, "externalIP", externalIP)
	return nil
}
