package models

type Node struct {
	NodeID     string `json:"node_id"`
	ObscuraKey string `json:"obscura_key"`
	ServerType string `json:"server_type"`
}
