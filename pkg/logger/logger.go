package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// 전역 로거 인스턴스
	globalLogger *zap.Logger
	// sugar 로거 인스턴스
	globalSugar *zap.SugaredLogger
)

// InitLogger는 로그 설정을 초기화하고 로거를 생성합니다
func InitLogger() (*zap.Logger, error) {
	// 현재 날짜로 로그 파일명 생성
	currentTime := time.Now()
	logFileName := fmt.Sprintf("logs/%s.log", currentTime.Format("2006-01-02"))

	// logs 디렉토리가 없으면 생성
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %v", err)
	}

	// 로그 설정
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		logFileName, // 파일 출력
		"stdout",    // 콘솔 출력
	}

	// 사람이 읽기 쉬운 시간 형식으로 변경
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 형식 사용 (2006-01-02T15:04:05.000Z)
	// 또는 아래처럼 커스텀 형식을 사용할 수도 있습니다
	// cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("로그 초기화 실패: %v", err)
	}

	// 글로벌 로거 설정
	globalLogger = logger
	globalSugar = logger.Sugar()

	return logger, nil
}

// GetLogger returns the global logger
func GetLogger() *zap.Logger {
	return globalLogger
}

// GetSugar returns the global sugared logger
func GetSugar() *zap.SugaredLogger {
	return globalSugar
}
