package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/pkg/models"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader   websocket.Upgrader
	store      func(*models.SystemMetrics) error
	cmdRepo    *repository.CommandRepository
	userRepo   *repository.UserRepository
	clients    sync.Map
	maxClients int
}

func NewServer(store func(*models.SystemMetrics) error, cmdRepo *repository.CommandRepository, userRepo *repository.UserRepository) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		store:      store,
		cmdRepo:    cmdRepo,
		userRepo:   userRepo,
		maxClients: 1000,
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

func (s *Server) handleMessage(conn *websocket.Conn, message []byte) {
	start := time.Now()
	var metrics models.SystemMetrics
	if err := json.Unmarshal(message, &metrics); err != nil {
		s.sendErrorResponse(conn, "메시지 파싱 오류")
		return
	}

	exists, err := s.userRepo.ExistsUserByObscuraKey(metrics.USER_ID)
	if err != nil || !exists {
		s.sendErrorResponse(conn, "사용자 조회 실패")
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

	// 응답 전송 전에 데드라인 설정
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 응답 전송
	response := models.WSResponse{
		Type:     "metrics_response",
		Commands: commands,
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("응답 전송 실패: %v", err)
	}

	elapsed := time.Since(start)
	log.Printf("응답 전송 완료: %v ms", elapsed.Milliseconds())
}

func (s *Server) sendErrorResponse(conn *websocket.Conn, errMsg string) {
	// 쓰기 작업 전에 데드라인 설정
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	response := models.WSResponse{
		Type:  "error",
		Error: errMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		log.Printf("에러 응답 전송 실패: %v", err)
	}
}

func (s *Server) handleDisconnect(clientID string, code int, text string) {
	log.Printf("클라이언트 %s 연결 종료 (코드: %d, 사유: %s)", clientID, code, text)
	if conn, ok := s.clients.LoadAndDelete(clientID); ok {
		if websocketConn, ok := conn.(*websocket.Conn); ok {
			websocketConn.Close()
		}
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	clientCount := 0
	s.clients.Range(func(_, _ interface{}) bool {
		clientCount++
		return true
	})
	if clientCount >= s.maxClients {
		log.Printf("최대 클라이언트 수 초과: %s", r.RemoteAddr)
		http.Error(w, "서버가 최대 용량에 도달했습니다", http.StatusServiceUnavailable)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("연결 업그레이드 실패: %v", err)
		return
	}

	// 읽기 데드라인만 설정하고 쓰기 데드라인은 각 쓰기 작업마다 설정하도록 수정
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	clientID := r.RemoteAddr
	s.clients.Store(clientID, conn)
	defer func() {
		s.clients.Delete(clientID)
		conn.Close()
	}()

	// ping-pong 처리를 위한 고루틴
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		// for range로 수정
		for range ticker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Ping 실패: %v", err)
				return
			}
		}
	}()

	log.Printf("새로운 클라이언트 연결: %s", clientID)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.handleDisconnect(clientID, websocket.CloseAbnormalClosure, err.Error())
			}
			break
		}

		// if len(message) > 1024*1024 {
		// 	s.sendErrorResponse(conn, "메시지가 너무 큽니다")
		// 	continue
		// }

		if messageType != websocket.TextMessage {
			s.sendErrorResponse(conn, "잘못된 메시지 타입")
			continue
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			done := make(chan struct{})
			go func() {
				s.handleMessage(conn, message)
				close(done)
			}()

			select {
			case <-ctx.Done():
				log.Printf("메시지 처리 시간 초과: %s", clientID)
				s.sendErrorResponse(conn, "처리 시간 초과")
			case <-done:
			}
		}()
	}
}
