package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
)

// CustomLogger는 기존 zap.SugaredLogger를 래핑하는 구조체
type CustomLogger struct {
	sugar *zap.SugaredLogger
	file  *os.File
}

var customLogger *CustomLogger

// InitCustomLogger는 커스텀 로거를 초기화합니다
func InitCustomLogger() error {
	// 일반 zap 로거 초기화
	zapLogger, err := InitLogger()
	if err != nil {
		return err
	}

	// 로그 파일 생성/오픈
	currentTime := time.Now()
	logFileName := fmt.Sprintf("logs/%s.log", currentTime.Format("2006-01-02"))

	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("로그 디렉토리 생성 실패: %v", err)
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("로그 파일 생성 실패: %v", err)
	}

	// 커스텀 로거 생성
	customLogger = &CustomLogger{
		sugar: zapLogger.Sugar(),
		file:  file,
	}

	return nil
}

// GetCustomLogger는 글로벌 커스텀 로거를 반환합니다
func GetCustomLogger() *CustomLogger {
	return customLogger
}

// 로그 레벨 메서드들 - 원래 zap.SugaredLogger와 동일한 인터페이스 제공

func (l *CustomLogger) Info(args ...interface{}) {
	l.log("INFO", args...)
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.logf("INFO", format, args...)
}

func (l *CustomLogger) Error(args ...interface{}) {
	l.log("ERROR", args...)
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.logf("ERROR", format, args...)
}

func (l *CustomLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logw("INFO", msg, keysAndValues...)
}

func (l *CustomLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logw("ERROR", msg, keysAndValues...)
}

func (l *CustomLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logw("WARN", msg, keysAndValues...)
}

func (l *CustomLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logw("DEBUG", msg, keysAndValues...)
}

func (l *CustomLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logw("FATAL", msg, keysAndValues...)
	os.Exit(1)
}

func (l *CustomLogger) Warn(args ...interface{}) {
	l.log("WARN", args...)
}

func (l *CustomLogger) Warnf(format string, args ...interface{}) {
	l.logf("WARN", format, args...)
}

func (l *CustomLogger) Debug(args ...interface{}) {
	l.log("DEBUG", args...)
}

func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.logf("DEBUG", format, args...)
}

func (l *CustomLogger) Fatal(args ...interface{}) {
	l.log("FATAL", args...)
	os.Exit(1)
}

func (l *CustomLogger) Fatalf(format string, args ...interface{}) {
	l.logf("FATAL", format, args...)
	os.Exit(1)
}

// 기타 필요한 메서드 추가...

// 내부 로깅 구현
func (l *CustomLogger) log(level string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 호출자 정보 가져오기
	_, file, line, _ := runtime.Caller(2) // 2단계 상위 호출자 가져오기

	// 메시지 생성
	msg := fmt.Sprint(args...)
	logLine := fmt.Sprintf("[%s] %s:%d: %s\n", timestamp, file, line, msg)

	// 파일과 콘솔에 출력
	fmt.Print(logLine)
	l.file.WriteString(logLine)
}

func (l *CustomLogger) logf(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 호출자 정보 가져오기
	_, file, line, _ := runtime.Caller(2) // 2단계 상위 호출자 가져오기

	// 메시지 생성
	msg := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] %s:%d: %s\n", timestamp, file, line, msg)

	// 파일과 콘솔에 출력
	fmt.Print(logLine)
	l.file.WriteString(logLine)
}

// 구조화된 로깅을 위한 logw 함수
func (l *CustomLogger) logw(level string, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 호출자 정보 가져오기
	_, file, line, _ := runtime.Caller(2)

	// 키-값 쌍 포맷팅
	var fields string
	for i := 0; i < len(keysAndValues); i += 2 {
		var value interface{}
		if i+1 < len(keysAndValues) {
			value = keysAndValues[i+1]
		} else {
			value = "<누락된 값>"
		}
		fields += fmt.Sprintf(" %v=%v", keysAndValues[i], value)
	}

	// 메시지 생성
	logLine := fmt.Sprintf("[%s] [%s] %s:%d: %s%s\n", timestamp, level, file, line, msg, fields)

	// 파일과 콘솔에 출력
	fmt.Print(logLine)
	l.file.WriteString(logLine)
}

// Close는 로거 리소스를 닫습니다
func (l *CustomLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
