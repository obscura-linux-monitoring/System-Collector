package main

import (
	config "system-collector/configs"
	"system-collector/internal/repository"
	"system-collector/internal/storage"
	"system-collector/internal/websocket"
	"system-collector/pkg/models"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"system-collector/pkg/logger"
)

func main() {
	// 로거 초기화
	if err := logger.InitCustomLogger(); err != nil {
		panic(fmt.Sprintf("로거 초기화 실패: %v", err))
	}

	sugar := logger.GetCustomLogger()
	defer sugar.Close()

	sugar.Infow("Starting System Collector")

	// 설정 로드
	if err := config.Load("configs/config.yaml"); err != nil {
		sugar.Errorw("설정 로드 실패", "error", err)
	}

	// 스토리지 초기화
	store, err := storage.NewInfluxDBClient()
	if err != nil {
		sugar.Errorw("Failed to initialize InfluxDB client", "error", err)
	}
	// defer store.Close()

	// PostgreSQL 클라이언트 초기화
	pgClient, err := storage.NewPostgresClient()
	if err != nil {
		sugar.Errorw("PostgreSQL 클라이언트 초기화 실패", "error", err)
	}
	// defer pgClient.Close()

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
					sugar.Errorw("메트릭스 저장 실패", "error", err)
				}
			}
		}()
	}

	// WebSocket 서버 초기화 (메트릭스 채널 전달)
	wsServer := websocket.NewServer(func(m *models.SystemMetrics) error {
		metricsChan <- m
		return nil
	}, cmdRepo, userRepo, nodeRepo)

	// 시그널 처리를 위한 채널 생성
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 별도의 고루틴에서 WebSocket 서버 시작
	go wsServer.Start()

	// 종료 시그널 대기
	sig := <-sigChan
	sugar.Infow("시스템 종료 신호 수신", "signal", sig.String())
	sugar.Infow("System Collector 종료 시작...")

	// 여기에서 명시적으로 종료 작업 수행
	sugar.Infow("데이터베이스 연결 종료 중...")
	pgClient.Close()
	sugar.Infow("스토리지 연결 종료 중...")
	store.Close()

	sugar.Infow("System Collector가 정상적으로 종료되었습니다")
}
