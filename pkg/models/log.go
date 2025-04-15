package models

import "time"

type LogMessage struct {
	NodeID    string    `json:"node_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Content   string    `json:"content"`
}
