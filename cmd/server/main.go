package main

import (
	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/internal/storage"
	"system-collector/internal/websocket"
	"system-collector/pkg/models"

	"fmt"

	"system-collector/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 로거 초기화
	logger, err := logger.InitLogger()
	if err != nil {
		panic(fmt.Sprintf("로거 초기화 실패: %v", err))
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Info("Starting System Collector")

	// 설정 로드
	if err := config.Load("configs/config.yaml"); err != nil {
		sugar.Error("설정 로드 실패", zap.Error(err))
	}

	// 스토리지 초기화
	store, err := storage.NewInfluxDBClient()
	if err != nil {
		sugar.Error("Failed to initialize InfluxDB client", zap.Error(err))
	}
	defer store.Close()

	// PostgreSQL 클라이언트 초기화
	pgClient, err := storage.NewPostgresClient()
	if err != nil {
		sugar.Error("PostgreSQL 클라이언트 초기화 실패", zap.Error(err))
	}
	defer pgClient.Close()

	// 커맨드 저장소 초기화
	cmdRepo := repository.NewCommandRepository(pgClient.GetDB())
	userRepo := repository.NewUserRepository(pgClient.GetDB())
	nodeRepo := repository.NewNodeRepository(pgClient.GetDB())
	// 버퍼가 있는 채널 생성
	metricsChan := make(chan *models.SystemMetrics, 1000)

	// 메트릭스 처리를 위한 워커 풀 생성
	for range [50]struct{}{} { // 워커 수는 필요에 따라 조정
		go func() {
			for metrics := range metricsChan {
				if err := store.StoreMetrics(metrics); err != nil {
					sugar.Error("메트릭스 저장 실패", zap.Error(err))
				}
			}
		}()
	}

	// WebSocket 서버 초기화 (메트릭스 채널 전달)
	wsServer := websocket.NewServer(func(m *models.SystemMetrics) error {
		metricsChan <- m
		return nil
	}, cmdRepo, userRepo, nodeRepo)

	// WebSocket 서버 시작
	wsServer.Start()
}
