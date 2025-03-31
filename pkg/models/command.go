package models

type Command struct {
	CommandID     int    `db:"command_id"`
	NodeID        string `db:"node_id"`
	CommandType   string `db:"command_type"`
	CommandStatus int16  `db:"command_status"`
}
