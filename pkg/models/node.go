package models

type Node struct {
	NodeID     string `json:"node_id"`
	ObscuraKey string `json:"obscura_key"`
	ServerType string `json:"server_type"`
	// ExternalIP는 노드의 외부 IP 주소를 나타냅니다
	ExternalIP string `json:"external_ip"`
}
