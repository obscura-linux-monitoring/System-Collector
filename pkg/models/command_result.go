package models

type CommandResult struct {
	CommandID     int    `json:"command_id"`
	NodeID        string `json:"node_id"`
	CommandType   string `json:"command_type"`
	CommandStatus int    `json:"command_status"`
	ResultStatus  int    `json:"result_status"`
	ResultMessage string `json:"result_message"`
	Target        string `json:"target"`
}
