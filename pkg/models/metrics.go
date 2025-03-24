// Package models는 시스템 메트릭스 수집을 위한 데이터 모델을 제공합니다.
package models

import (
	"encoding/json"
	"log"
	"time"
)

// SystemMetrics는 시스템의 전반적인 상태 정보를 포함하는 구조체입니다.
// 시스템의 주요 구성 요소인 CPU, 메모리, 디스크, 네트워크 사용량과
// 실행 중인 프로세스 및 도커 컨테이너 정보를 포함합니다.
type SystemMetrics struct {
	// Key는 메트릭스의 고유 식별자입니다
	Key string `json:"key"`
	// Timestamp는 메트릭스가 수집된 시간을 나타냅니다
	Timestamp time.Time `json:"timestamp"`
	// CPU는 프로세서 관련 메트릭스를 포함합니다
	CPU CPUMetrics `json:"cpu"`
	// Memory는 메모리 사용량 관련 메트릭스를 포함합니다
	Memory MemoryMetrics `json:"memory"`
	// Disk는 디스크별 사용량 정보를 포함합니다
	Disk []DiskMetrics `json:"disk,omitempty"`
	// Network는 네트워크 인터페이스별 통계를 포함합니다
	Network []NetworkMetrics `json:"network,omitempty"`
	// Processes는 실행 중인 프로세스 목록을 포함합니다
	Processes []ProcessInfo `json:"processes,omitempty"`
	// Containers는 실행 중인 도커 컨테이너 목록을 포함합니다
	Containers []DockerContainer `json:"containers,omitempty"`
}

// CPUMetrics는 CPU 관련 메트릭스를 포함하는 구조체입니다.
type CPUMetrics struct {
	// UsagePercent는 전체 CPU 사용률을 백분율로 나타냅니다
	UsagePercent float64 `json:"usage_percent"`
	// PerCorePercent는 코어별 CPU 사용률을 백분율로 나타냅니다
	PerCorePercent map[string]float64 `json:"per_core_percent"`
	// Temperature는 CPU 온도를 섭씨로 나타냅니다
	Temperature float64 `json:"temperature"`
	// LoadAvg1은 1분 평균 시스템 로드를 나타냅니다
	LoadAvg1 float64 `json:"load_avg_1"`
	// LoadAvg5는 5분 평균 시스템 로드를 나타냅니다
	LoadAvg5 float64 `json:"load_avg_5"`
	// LoadAvg15는 15분 평균 시스템 로드를 나타냅니다
	LoadAvg15 float64 `json:"load_avg_15"`
}

// MemoryMetrics는 시스템 메모리 사용량 정보를 포함하는 구조체입니다.
type MemoryMetrics struct {
	// Total은 전체 메모리 크기를 바이트 단위로 나타냅니다
	Total int64 `json:"total"`
	// Used는 사용 중인 메모리를 바이트 단위로 나타냅니다
	Used int64 `json:"used"`
	// Free는 사용 가능한 메모리를 바이트 단위로 나타냅니다
	Free int64 `json:"free"`
	// UsagePercent는 메모리 사용률을 백분율로 나타냅니다
	UsagePercent float64 `json:"usage_percent"`
}

// DiskMetrics는 개별 디스크의 사용량 정보를 포함하는 구조체입니다.
type DiskMetrics struct {
	// Device는 디스크 장치명을 나타냅니다
	Device string `json:"device"`
	// MountPoint는 디스크의 마운트 위치를 나타냅니다
	MountPoint string `json:"mount_point"`
	// Total은 전체 디스크 용량을 바이트 단위로 나타냅니다
	Total int64 `json:"total"`
	// Used는 사용 중인 디스크 공간을 바이트 단위로 나타냅니다
	Used int64 `json:"used"`
	// Free는 사용 가능한 디스크 공간을 바이트 단위로 나타냅니다
	Free int64 `json:"free"`
	// UsagePercent는 디스크 사용률을 백분율로 나타냅니다
	UsagePercent float64 `json:"usage_percent"`
}

// NetworkMetrics는 네트워크 인터페이스의 통계 정보를 포함하는 구조체입니다.
type NetworkMetrics struct {
	// Interface는 네트워크 인터페이스 이름을 나타냅니다
	Interface string `json:"interface"`
	// RxBytes는 수신한 총 바이트를 나타냅니다
	RxBytes int64 `json:"rx_bytes"`
	// TxBytes는 전송한 총 바이트를 나타냅니다
	TxBytes int64 `json:"tx_bytes"`
	// RxPackets는 수신한 총 패킷 수를 나타냅니다
	RxPackets int64 `json:"rx_packets"`
	// TxPackets는 전송한 총 패킷 수를 나타냅니다
	TxPackets int64 `json:"tx_packets"`
	// RxErrors는 수신 중 발생한 오류 수를 나타냅니다
	RxErrors int64 `json:"rx_errors"`
	// TxErrors는 전송 중 발생한 오류 수를 나타냅니다
	TxErrors int64 `json:"tx_errors"`
	// RxDropped는 수신 중 드롭된 패킷 수를 나타냅니다
	RxDropped int64 `json:"rx_dropped"`
	// TxDropped는 전송 중 드롭된 패킷 수를 나타냅니다
	TxDropped int64 `json:"tx_dropped"`
}

// ProcessInfo는 개별 프로세스의 상태 정보를 포함하는 구조체입니다.
type ProcessInfo struct {
	// PID는 프로세스 식별자를 나타냅니다
	PID int `json:"pid"`
	// Name은 프로세스 이름을 나타냅니다
	Name string `json:"name"`
	// CPU는 프로세스의 CPU 사용률을 백분율로 나타냅니다
	CPU float64 `json:"cpu"`
	// Memory는 프로세스의 메모리 사용률을 백분율로 나타냅니다
	Memory float64 `json:"memory"`
	// User는 프로세스를 실행한 사용자를 나타냅니다
	User string `json:"user"`
	// Command는 프로세스 실행 명령어를 나타냅니다
	Command string `json:"command"`
	// ThreadCount는 프로세스의 스레드 수를 나타냅니다
	ThreadCount int `json:"thread_count"`
	// Status는 프로세스의 현재 상태를 나타냅니다
	Status string `json:"status"`
}

// DockerContainer는 도커 컨테이너의 상태 정보를 포함하는 구조체입니다.
type DockerContainer struct {
	// ID는 컨테이너의 고유 식별자를 나타냅니다
	ID string `json:"id"`
	// Name은 컨테이너의 이름을 나타냅니다
	Name string `json:"name"`
	// Image는 컨테이너의 이미지 정보를 나타냅니다
	Image string `json:"image"`
	// Status는 컨테이너의 현재 상태를 나타냅니다
	Status string `json:"status"`
	// CPU는 컨테이너의 CPU 사용률을 백분율로 나타냅니다
	CPU float64 `json:"cpu"`
	// Memory는 컨테이너의 메모리 사용률을 백분율로 나타냅니다
	Memory float64 `json:"memory"`
	// NetworkRx는 컨테이너가 수신한 네트워크 데이터량을 바이트 단위로 나타냅니다
	NetworkRx int64 `json:"network_rx"`
	// NetworkTx는 컨테이너가 전송한 네트워크 데이터량을 바이트 단위로 나타냅니다
	NetworkTx int64 `json:"network_tx"`
	// BlockRead는 컨테이너가 읽은 블록 디바이스 데이터량을 바이트 단위로 나타냅니다
	BlockRead int64 `json:"block_read"`
	// BlockWrite는 컨테이너가 쓴 블록 디바이스 데이터량을 바이트 단위로 나타냅니다
	BlockWrite int64 `json:"block_write"`
	// Labels는 컨테이너에 설정된 레이블 목록을 나타냅니다
	Labels []struct {
		Key   string `json:"label_key"`
		Value string `json:"label_value"`
	} `json:"labels,omitempty"`
}

// ToJSON은 주어진 인터페이스를 JSON 바이트 배열로 변환합니다.
// 변환 중 오류가 발생하면 빈 JSON 객체({})를 반환합니다.
//
// 매개변수:
//   - v: JSON으로 변환할 인터페이스
//
// 반환값:
//   - []byte: JSON으로 변환된 바이트 배열
func ToJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return []byte("{}")
	}
	return data
}
