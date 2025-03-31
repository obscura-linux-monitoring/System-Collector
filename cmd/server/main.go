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
	userRepo := repository.NewUserRepository(pgClient.GetDB())

	// 버퍼가 있는 채널 생성
	metricsChan := make(chan *models.SystemMetrics, 1000)

	// 메트릭스 처리를 위한 워커 풀 생성
	for range [50]struct{}{} { // 워커 수는 필요에 따라 조정
		go func() {
			for metrics := range metricsChan {
				if err := store.StoreMetrics(metrics); err != nil {
					log.Printf("메트릭스 저장 실패: %v", err)
				}
			}
		}()
	}

	// WebSocket 서버 초기화 (메트릭스 채널 전달)
	wsServer := websocket.NewServer(func(m *models.SystemMetrics) error {
		metricsChan <- m
		return nil
	}, cmdRepo, userRepo)

	// WebSocket 서버 시작
	wsServer.Start()
}
