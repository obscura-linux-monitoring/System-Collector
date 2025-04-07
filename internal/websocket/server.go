package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/pkg/logger"
	"system-collector/pkg/models"

	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader   websocket.Upgrader
	store      func(*models.SystemMetrics) error
	cmdRepo    *repository.CommandRepository
	userRepo   *repository.UserRepository
	nodeRepo   *repository.NodeRepository
	nodeList   []*models.Node
	clients    sync.Map
	maxClients int
}

func NewServer(store func(*models.SystemMetrics) error, cmdRepo *repository.CommandRepository, userRepo *repository.UserRepository, nodeRepo *repository.NodeRepository) *Server {
	sugar := logger.GetCustomLogger()
	sugar.Infow("Server 초기화 중")

	// 사용자 목록 조회
	nodes, err := nodeRepo.GetAllNodes()
	if err != nil {
		sugar.Errorw("사용자 목록 조회 실패", "error", err)
		nodes = []*models.Node{} // 빈 배열로 초기화
	}

	return &Server{
		upgrader: websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		store:      store,
		cmdRepo:    cmdRepo,
		userRepo:   userRepo,
		nodeRepo:   nodeRepo,
		nodeList:   nodes, // 조회한 사용자 목록 설정
		maxClients: 1000,
	}
}

func (s *Server) Start() {
	sugar := logger.GetCustomLogger()
	sugar.Infow("WebSocket server 시작 중")

	http.HandleFunc("/ws", s.handleConnections)

	cfg := config.Get()
	sugar.Infow("WebSocket server starting", "port", cfg.Server.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		sugar.Errorf("WebSocket server failed to start: %v", err)
	}
}

func (s *Server) handleMessage(conn *websocket.Conn, message []byte) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("handleMessage 시작")

	start := time.Now()
	var metrics models.SystemMetrics
	if err := json.Unmarshal(message, &metrics); err != nil {
		sugar.Errorw("메시지 파싱 오류", "error", err)
		s.sendErrorResponse(conn, "메시지 파싱 오류")
		return
	}

	sugar.Infof("커맨드 결과 : %v", metrics.CommandResults)

	if len(metrics.CommandResults) > 0 {
		sugar.Infof("커맨드 결과 전송 시작: %d개", len(metrics.CommandResults))

		// 커맨드 결과 JSON 변환
		commandResultJSON, err := json.Marshal(metrics.CommandResults)
		if err != nil {
			sugar.Errorw("커맨드 결과 JSON 변환 실패", "error", err)
		} else {
			s.sendCommandResults(commandResultJSON, metrics.Key, metrics.USER_ID)
		}
	}

	exists, err := s.userRepo.ExistsUserByObscuraKey(metrics.USER_ID)
	if err != nil || !exists {
		sugar.Errorw("사용자 조회 실패", "error", err)
		s.sendErrorResponse(conn, "사용자 조회 실패")
		return
	}

	// Key 검증
	if metrics.Key == "" {
		sugar.Errorw("키가 없는 메트릭스")
		s.sendErrorResponse(conn, "키가 없는 메트릭스")
		return
	}

	nodeExists := false
	for _, node := range s.nodeList {
		if node.NodeID == metrics.Key {
			nodeExists = true
			break
		}
	}

	if !nodeExists {
		node := models.Node{
			NodeID:     metrics.Key,
			ObscuraKey: metrics.USER_ID,
			ServerType: "debug", // TODO: 서버 타입 추가
		}
		err := s.nodeRepo.CreateNode(&node)
		if err != nil {
			sugar.Errorw("노드 생성 실패", "error", err)
		} else {
			sugar.Infof("노드 생성 성공: %s", metrics.Key)
			s.nodeList = append(s.nodeList, &models.Node{
				NodeID: metrics.Key,
			})
		}
	}

	sugar.Infof("메트릭스 키: %s", metrics.Key)

	// 메트릭스 저장
	if err := s.store(&metrics); err != nil {
		sugar.Errorw("메트릭스 저장 실패", "error", err)
		s.sendErrorResponse(conn, "메트릭스 저장 실패")
		return
	}

	// 해당 노드의 명령어 조회
	commands, err := s.cmdRepo.GetCommandsByNodeID(metrics.Key)
	if err != nil {
		sugar.Errorw("명령어 조회 실패", "error", err)
		// 명령어 조회 실패해도 메트릭스는 정상 응답
		commands = []models.Command{}
	} else {
		sugar.Infof("명령어 조회 성공: %v", commands)
		if len(commands) > 0 {
			err = s.cmdRepo.DeleteCommandsByNodeID(metrics.Key)
			if err != nil {
				sugar.Errorw("명령어 삭제 실패", "error", err)
			}
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
		sugar.Errorw("응답 전송 실패", "error", err)
	}

	elapsed := time.Since(start)
	sugar.Infof("응답 전송 완료: %v ms", elapsed.Milliseconds())
}

// sendCommandResults는 명령어 실행 결과를 REST API로 전송하는 함수입니다
func (s *Server) sendCommandResults(commandResultJSON []byte, nodeID, userID string) {
	sugar := logger.GetCustomLogger()

	// HTTP 클라이언트 생성
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// REST API 요청 생성
	req, err := http.NewRequest("POST", "http://1.209.148.143:8000/api/command-results", bytes.NewBuffer(commandResultJSON))
	if err != nil {
		sugar.Errorw("커맨드 결과 요청 생성 실패", "error", err)
		return
	}

	// 요청 헤더 설정
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Node-ID", nodeID)
	req.Header.Set("X-User-ID", userID)

	// 요청 전송
	resp, err := client.Do(req)
	if err != nil {
		sugar.Errorw("커맨드 결과 전송 실패", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		sugar.Infof("커맨드 결과 전송 성공: %d", resp.StatusCode)
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		sugar.Errorw("커맨드 결과 전송 실패", "status", resp.StatusCode, "response", string(bodyBytes))
	}
}

func (s *Server) sendErrorResponse(conn *websocket.Conn, errMsg string) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("sendErrorResponse 시작")

	// 쓰기 작업 전에 데드라인 설정
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	response := models.WSResponse{
		Type:  "error",
		Error: errMsg,
	}
	if err := conn.WriteJSON(response); err != nil {
		sugar.Errorw("에러 응답 전송 실패", "error", err)
	}
}

func (s *Server) handleDisconnect(clientID string, code int, text string) {
	sugar := logger.GetCustomLogger()
	sugar.Infof("클라이언트 %s 연결 종료 (코드: %d, 사유: %s)", clientID, code, text)
	if conn, ok := s.clients.LoadAndDelete(clientID); ok {
		if websocketConn, ok := conn.(*websocket.Conn); ok {
			websocketConn.Close()
		}
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("handleConnections 시작")

	clientCount := 0
	s.clients.Range(func(_, _ interface{}) bool {
		clientCount++
		return true
	})
	if clientCount >= s.maxClients {
		sugar.Infof("최대 클라이언트 수 초과: %s", r.RemoteAddr)
		http.Error(w, "서버가 최대 용량에 도달했습니다", http.StatusServiceUnavailable)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		sugar.Errorw("연결 업그레이드 실패", "error", err)
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
				sugar.Errorw("Ping 실패", "error", err)
				return
			}
		}
	}()

	sugar.Infof("새로운 클라이언트 연결: %s", clientID)

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
				sugar.Infof("메시지 처리 시간 초과: %s", clientID)
				s.sendErrorResponse(conn, "처리 시간 초과")
			case <-done:
			}
		}()
	}
}
