package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	config "system-collector/configs"
	"system-collector/pkg/models"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	store    func(*models.SystemMetrics) error
}

func NewServer(store func(*models.SystemMetrics) error) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 개발용, 실제 환경에서는 적절한 검사 필요
			},
		},
		store: store,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.handleConnections)

	cfg := config.Get()
	log.Printf("WebSocket server starting on :%d", cfg.Server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		log.Fatalf("WebSocket server failed to start: %v", err)
	}
}

// 별도의 함수로 분리하여 재사용성 높임
func handleDisconnect(r *http.Request, code int, text string) {
	log.Printf("Client %s disconnected with code %d: %s", r.RemoteAddr, code, text)
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// SetCloseHandler로 disconnect 처리
	conn.SetCloseHandler(func(code int, text string) error {
		handleDisconnect(r, code, text)
		return nil
	})

	// 클라이언트 접속 로그 추가
	log.Printf("New client connected from %s", r.RemoteAddr)
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error from %s: %v", r.RemoteAddr, err)
			}
			handleDisconnect(r, websocket.CloseAbnormalClosure, err.Error())
			break
		}

		// 모든 수신 메시지에 대한 로깅 추가
		log.Printf("Message received from %s (length: %d bytes)", r.RemoteAddr, len(message))

		// 메시지 처리를 별도 고루틴으로 분리하여 다음 메시지를 즉시 수신할 수 있도록 함
		go func(msg []byte) {
			// 수신된 메시지를 SystemMetrics로 파싱
			var metrics models.SystemMetrics
			if err := json.Unmarshal(msg, &metrics); err != nil {
				log.Printf("Error parsing message: %v", err)
				return
			}

			// Key 검증
			if metrics.Key == "" {
				log.Printf("Received metrics without key")
				return
			}

			// 받은 메트릭 데이터를 보기 좋게 출력
			log.Printf("Successfully processed metrics from %s with timestamp %v", r.RemoteAddr, metrics.Timestamp)

			// InfluxDB에 저장
			if err := s.store(&metrics); err != nil {
				log.Printf("Error storing metrics: %v", err)
				return
			}
			log.Printf("Metrics successfully stored in database")
		}(message)
	}
}
