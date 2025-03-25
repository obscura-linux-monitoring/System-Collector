package main

import (
	"log"
	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/internal/storage"
	"system-collector/internal/websocket"
	"system-collector/pkg/models"
)

func main() {
	log.Println("Starting System Collector")

	// 설정 로드
	if err := config.Load("configs/config.yaml"); err != nil {
		log.Fatalf("설정 로드 실패: %v", err)
	}

	// 스토리지 초기화
	store, err := storage.NewInfluxDBClient()
	if err != nil {
		log.Fatalf("Failed to initialize InfluxDB client: %v", err)
	}
	defer store.Close()

	// PostgreSQL 클라이언트 초기화
	pgClient, err := storage.NewPostgresClient()
	if err != nil {
		log.Fatalf("PostgreSQL 클라이언트 초기화 실패: %v", err)
	}
	defer pgClient.Close()

	// 커맨드 저장소 초기화
	cmdRepo := repository.NewCommandRepository(pgClient.GetDB())

	// WebSocket 서버 초기화 (cmdRepo 추가)
	wsServer := websocket.NewServer(func(m *models.SystemMetrics) error {
		return store.StoreMetrics(m)
	}, cmdRepo)

	// WebSocket 서버 시작
	wsServer.Start()
}
