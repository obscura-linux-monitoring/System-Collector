package models

// WSResponse는 WebSocket 응답 구조체입니다
type WSResponse struct {
	Type   string `json:"type"`
	Result string `json:"result"`
}
