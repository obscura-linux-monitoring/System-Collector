package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/pkg/models"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader websocket.Upgrader
	store    func(*models.SystemMetrics) error
	cmdRepo  *repository.CommandRepository
}

func NewServer(store func(*models.SystemMetrics) error, cmdRepo *repository.CommandRepository) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 개발용, 실제 환경에서는 적절한 검사 필요
			},
		},
		store:   store,
		cmdRepo: cmdRepo,
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

func (s *Server) handleMessage(conn *websocket.Conn, message []byte) {
	var metrics models.SystemMetrics
	if err := json.Unmarshal(message, &metrics); err != nil {
		s.sendErrorResponse(conn, "메시지 파싱 오류")
		return
	}

	// Key 검증
	if metrics.Key == "" {
		s.sendErrorResponse(conn, "키가 없는 메트릭스")
		return
	}

	log.Printf("메트릭스 키: %s", metrics.Key)

	// 메트릭스 저장
	if err := s.store(&metrics); err != nil {
		s.sendErrorResponse(conn, "메트릭스 저장 실패")
		return
	}

	// 해당 노드의 명령어 조회
	commands, err := s.cmdRepo.GetCommandsByNodeID(metrics.Key)
	if err != nil {
		log.Printf("명령어 조회 실패: %v", err)
		// 명령어 조회 실패해도 메트릭스는 정상 응답
		commands = []models.Command{}
	} else {
		log.Printf("명령어 조회 성공: %v", commands)
		err = s.cmdRepo.DeleteCommandsByNodeID(metrics.Key)
		if err != nil {
			log.Printf("명령어 삭제 실패: %v", err)
		}
	}

	// 응답 전송
	response := models.WSResponse{
		Type:     "metrics_response",
		Commands: commands,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("응답 전송 실패: %v", err)
	}
}

func (s *Server) sendErrorResponse(conn *websocket.Conn, errMsg string) {
	response := models.WSResponse{
		Type:  "error",
		Error: errMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("에러 응답 전송 실패: %v", err)
	}
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

		// 메시지 처리를 고루틴으로 실행
		go s.handleMessage(conn, message)
	}
}
