package main

import (
	"log"
	config "system-collector/configs"
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

	// WebSocket 서버 초기화
	wsServer := websocket.NewServer(func(m *models.SystemMetrics) error {
		return store.StoreMetrics(m)
	})

	// WebSocket 서버 시작
	wsServer.Start()
}
