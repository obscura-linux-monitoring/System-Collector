package models

type Node struct {
	NodeID     string `json:"node_id"`
	ObscuraKey string `json:"obscura_key"`
	ServerType bool   `json:"server_type"`
}
