// Package models는 시스템 메트릭스 수집을 위한 데이터 모델을 제공합니다.
package models

import (
	"encoding/json"
	"log"
	"time"
)

// SystemMetrics는 시스템의 전반적인 상태 정보를 포함하는 구조체입니다.
type SystemMetrics struct {
	USER_ID string `json:"user_id"`
	// Key는 메트릭스의 고유 식별자입니다
	Key string `json:"key"`
	// CommandResults는 명령어 실행 결과를 포함합니다
	CommandResults []CommandResult `json:"command_results"`
	// Timestamp는 메트릭스가 수집된 시간을 나타냅니다
	Timestamp time.Time `json:"timestamp"`
	// System은 시스템 기본 정보를 포함합니다
	System SystemInfo `json:"system"`
	// CPU는 프로세서 관련 메트릭스를 포함합니다
	CPU CPUMetrics `json:"cpu"`
	// Memory는 메모리 사용량 관련 메트릭스를 포함합니다
	Memory MemoryMetrics `json:"memory"`
	// Disk는 디스크별 사용량 정보를 포함합니다
	Disk []DiskMetrics `json:"disk"`
	// Network는 네트워크 인터페이스별 통계를 포함합니다
	Network []NetworkMetrics `json:"network"`
	// Processes는 실행 중인 프로세스 목록을 포함합니다
	Processes []ProcessInfo `json:"processes"`
	// Containers는 실행 중인 도커 컨테이너 목록을 포함합니다
	Containers []DockerContainer `json:"containers"`
	// Services는 시스템 서비스 정보를 포함합니다
	Services []ServiceInfo `json:"services"`
}

// SystemInfo는 시스템 기본 정보를 포함하는 구조체입니다.
type SystemInfo struct {
	// Hostname은 시스템의 호스트명을 나타냅니다
	Hostname string `json:"hostname"`
	// OSName은 운영체제 이름을 나타냅니다
	OSName string `json:"os_name"`
	// OSVersion은 운영체제 버전을 나타냅니다
	OSVersion string `json:"os_version"`
	// OSArchitecture는 시스템 아키텍처를 나타냅니다
	OSArchitecture string `json:"os_architecture"`
	// OSKernelVersion은 커널 버전을 나타냅니다
	OSKernelVersion string `json:"os_kernel_version"`
	// BootTime은 시스템 부팅 시간을 나타냅니다
	BootTime int64 `json:"boot_time"`
	// Uptime은 시스템 가동 시간을 나타냅니다
	Uptime int64 `json:"uptime"`
}

// CPUMetrics는 CPU 관련 메트릭스를 포함하는 구조체입니다.
type CPUMetrics struct {
	// Architecture는 CPU 아키텍처를 나타냅니다
	Architecture string `json:"architecture"`
	// Model은 CPU 모델명을 나타냅니다
	Model string `json:"model"`
	// Vendor는 CPU 제조사를 나타냅니다
	Vendor string `json:"vendor"`
	// CacheSize는 CPU 캐시 크기를 나타냅니다
	CacheSize int `json:"cache_size"`
	// ClockSpeed는 CPU 클럭 속도를 나타냅니다
	ClockSpeed float64 `json:"clock_speed"`
	// TotalCores는 물리 코어 수를 나타냅니다
	TotalCores int `json:"total_cores"`
	// TotalLogicalCores는 논리 코어 수를 나타냅니다
	TotalLogicalCores int `json:"total_logical_cores"`
	// Usage는 전체 CPU 사용률을 나타냅니다
	Usage float64 `json:"usage"`
	// Temperature는 CPU 온도를 나타냅니다
	Temperature float64 `json:"temperature"`
	// HasVMX는 가상화 지원 여부를 나타냅니다
	HasVMX bool `json:"has_vmx"`
	// HasSVM은 AMD-V 지원 여부를 나타냅니다
	HasSVM bool `json:"has_svm"`
	// HasAVX는 AVX 지원 여부를 나타냅니다
	HasAVX bool `json:"has_avx"`
	// HasAVX2는 AVX2 지원 여부를 나타냅니다
	HasAVX2 bool `json:"has_avx2"`
	// HasNEON은 NEON 지원 여부를 나타냅니다
	HasNEON bool `json:"has_neon"`
	// HasSVE는 SVE 지원 여부를 나타냅니다
	HasSVE bool `json:"has_sve"`
	// IsHyperthreading은 하이퍼스레딩 지원 여부를 나타냅니다
	IsHyperthreading bool `json:"is_hyperthreading"`
	// Cores는 코어별 상세 정보를 나타냅니다
	Cores []CPUCore `json:"cores"`
}

// CPUCore는 개별 CPU 코어의 정보를 포함하는 구조체입니다.
type CPUCore struct {
	// ID는 코어 ID를 나타냅니다
	ID int `json:"id"`
	// Usage는 코어별 사용률을 나타냅니다
	Usage float64 `json:"usage"`
	// Temperature는 코어별 온도를 나타냅니다
	Temperature float64 `json:"temperature"`
}

// MemoryMetrics는 시스템 메모리 사용량 정보를 포함하는 구조체입니다.
type MemoryMetrics struct {
	// Total은 전체 메모리 크기를 바이트 단위로 나타냅니다
	Total int64 `json:"total"`
	// Used는 사용 중인 메모리를 바이트 단위로 나타냅니다
	Used int64 `json:"used"`
	// Free는 사용 가능한 메모리를 바이트 단위로 나타냅니다
	Free int64 `json:"free"`
	// Available은 실제 사용 가능한 메모리를 바이트 단위로 나타냅니다
	Available int64 `json:"available"`
	// Buffers는 버퍼로 사용 중인 메모리를 바이트 단위로 나타냅니다
	Buffers int64 `json:"buffers"`
	// Cached는 캐시로 사용 중인 메모리를 바이트 단위로 나타냅니다
	Cached int64 `json:"cached"`
	// SwapTotal은 전체 스왑 크기를 바이트 단위로 나타냅니다
	SwapTotal int64 `json:"swap_total"`
	// SwapUsed는 사용 중인 스왑을 바이트 단위로 나타냅니다
	SwapUsed int64 `json:"swap_used"`
	// SwapFree는 사용 가능한 스왑을 바이트 단위로 나타냅니다
	SwapFree int64 `json:"swap_free"`
	// UsagePercent는 메모리 사용률을 백분율로 나타냅니다
	UsagePercent float64 `json:"usage_percent"`
}

// DiskMetrics는 개별 디스크의 사용량 정보를 포함하는 구조체입니다.
type DiskMetrics struct {
	// Device는 디스크 장치명을 나타냅니다
	Device string `json:"device"`
	// MountPoint는 디스크의 마운트 위치를 나타냅니다
	MountPoint string `json:"mount_point"`
	// FilesystemType은 파일시스템 종류를 나타냅니다
	FilesystemType string `json:"filesystem_type"`
	// Total은 전체 디스크 용량을 바이트 단위로 나타냅니다
	Total int64 `json:"total"`
	// Used는 사용 중인 디스크 공간을 바이트 단위로 나타냅니다
	Used int64 `json:"used"`
	// Free는 사용 가능한 디스크 공간을 바이트 단위로 나타냅니다
	Free int64 `json:"free"`
	// InodesTotal은 전체 아이노드 수를 나타냅니다
	InodesTotal int64 `json:"inodes_total"`
	// InodesUsed는 사용 중인 아이노드 수를 나타냅니다
	InodesUsed int64 `json:"inodes_used"`
	// InodesFree는 사용 가능한 아이노드 수를 나타냅니다
	InodesFree int64 `json:"inodes_free"`
	// UsagePercent는 디스크 사용률을 백분율로 나타냅니다
	UsagePercent float64 `json:"usage_percent"`
	// ErrorFlag는 디스크 오류 여부를 나타냅니다
	ErrorFlag bool `json:"error_flag"`
	// ErrorMessage는 디스크 오류 메시지를 나타냅니다
	ErrorMessage string `json:"error_message"`
	// IOStats는 디스크 I/O 통계를 나타냅니다
	IOStats DiskIOStats `json:"io_stats"`
}

// DiskIOStats는 디스크 I/O 통계 정보를 포함하는 구조체입니다.
type DiskIOStats struct {
	// ReadBytes는 읽은 총 바이트를 나타냅니다
	ReadBytes int64 `json:"read_bytes"`
	// WriteBytes는 쓴 총 바이트를 나타냅니다
	WriteBytes int64 `json:"write_bytes"`
	// Reads는 읽기 작업 횟수를 나타냅니다
	Reads int64 `json:"reads"`
	// Writes는 쓰기 작업 횟수를 나타냅니다
	Writes int64 `json:"writes"`
	// ReadBytesPerSec는 초당 읽은 바이트를 나타냅니다
	ReadBytesPerSec float64 `json:"read_bytes_per_sec"`
	// WriteBytesPerSec는 초당 쓴 바이트를 나타냅니다
	WriteBytesPerSec float64 `json:"write_bytes_per_sec"`
	// ReadsPerSec는 초당 읽기 작업 횟수를 나타냅니다
	ReadsPerSec float64 `json:"reads_per_sec"`
	// WritesPerSec는 초당 쓰기 작업 횟수를 나타냅니다
	WritesPerSec float64 `json:"writes_per_sec"`
	// IOInProgress는 진행 중인 I/O 작업 수를 나타냅니다
	IOInProgress int64 `json:"io_in_progress"`
	// IOTime은 I/O 작업에 소요된 시간을 나타냅니다
	IOTime int64 `json:"io_time"`
	// ReadTime은 읽기 작업에 소요된 시간을 나타냅니다
	ReadTime int64 `json:"read_time"`
	// WriteTime은 쓰기 작업에 소요된 시간을 나타냅니다
	WriteTime int64 `json:"write_time"`
	// ErrorFlag는 I/O 오류 여부를 나타냅니다
	ErrorFlag bool `json:"error_flag"`
}

// NetworkMetrics는 네트워크 인터페이스의 통계 정보를 포함하는 구조체입니다.
type NetworkMetrics struct {
	// Interface는 네트워크 인터페이스 이름을 나타냅니다
	Interface string `json:"interface"`
	// IP는 인터페이스의 IP 주소를 나타냅니다
	IP string `json:"ip"`
	// MAC은 인터페이스의 MAC 주소를 나타냅니다
	MAC string `json:"mac"`
	// MTU는 최대 전송 단위를 나타냅니다
	MTU int `json:"mtu"`
	// Speed는 인터페이스 속도를 나타냅니다
	Speed uint64 `json:"speed"`
	// Status는 인터페이스 상태를 나타냅니다
	Status string `json:"status"`
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
	// RxBytesPerSec는 초당 수신 바이트를 나타냅니다
	RxBytesPerSec float64 `json:"rx_bytes_per_sec"`
	// TxBytesPerSec는 초당 전송 바이트를 나타냅니다
	TxBytesPerSec float64 `json:"tx_bytes_per_sec"`
}

// ProcessInfo는 개별 프로세스의 상태 정보를 포함하는 구조체입니다.
type ProcessInfo struct {
	// PID는 프로세스 식별자를 나타냅니다
	PID int `json:"pid"`
	// PPID는 부모 프로세스 식별자를 나타냅니다
	PPID int `json:"ppid"`
	// Name은 프로세스 이름을 나타냅니다
	Name string `json:"name"`
	// Command는 프로세스 실행 명령어를 나타냅니다
	Command string `json:"command"`
	// User는 프로세스를 실행한 사용자를 나타냅니다
	User string `json:"user"`
	// Status는 프로세스의 현재 상태를 나타냅니다
	Status string `json:"status"`
	// CPUTime은 CPU 사용 시간을 나타냅니다
	CPUTime float64 `json:"cpu_time"`
	// CPUUsage는 CPU 사용률을 나타냅니다
	CPUUsage float64 `json:"cpu_usage"`
	// MemoryRSS는 실제 메모리 사용량을 나타냅니다
	MemoryRSS int64 `json:"memory_rss"`
	// MemoryVSZ는 가상 메모리 사용량을 나타냅니다
	MemoryVSZ int64 `json:"memory_vsz"`
	// Nice는 프로세스 우선순위를 나타냅니다
	Nice int `json:"nice"`
	// Threads는 스레드 수를 나타냅니다
	Threads int `json:"threads"`
	// OpenFiles는 열린 파일 수를 나타냅니다
	OpenFiles int `json:"open_files"`
	// StartTime은 프로세스 시작 시간을 나타냅니다
	StartTime int64 `json:"start_time"`
	// IOReadBytes는 읽은 총 바이트를 나타냅니다
	IOReadBytes int64 `json:"io_read_bytes"`
	// IOWriteBytes는 쓴 총 바이트를 나타냅니다
	IOWriteBytes int64 `json:"io_write_bytes"`
}

// DockerContainer는 도커 컨테이너의 상태 정보를 포함하는 구조체입니다.
type DockerContainer struct {
	// ID는 컨테이너의 고유 식별자를 나타냅니다
	ID string `json:"container_id"`
	// Name은 컨테이너의 이름을 나타냅니다
	Name string `json:"container_name"`
	// Image는 컨테이너의 이미지 정보를 나타냅니다
	Image string `json:"container_image"`
	// Status는 컨테이너의 현재 상태를 나타냅니다
	Status string `json:"container_status"`
	// Created는 컨테이너 생성 시간을 나타냅니다
	Created string `json:"container_created"`
	// CPUUsage는 CPU 사용률을 나타냅니다
	CPUUsage float64 `json:"cpu_usage"`
	// MemoryUsage는 메모리 사용량을 나타냅니다
	MemoryUsage int64 `json:"memory_usage"`
	// MemoryLimit는 메모리 제한을 나타냅니다
	MemoryLimit int64 `json:"memory_limit"`
	// MemoryPercent는 메모리 사용률을 나타냅니다
	MemoryPercent float64 `json:"memory_percent"`
	// NetworkRxBytes는 수신한 총 바이트를 나타냅니다
	NetworkRxBytes int64 `json:"network_rx_bytes"`
	// NetworkTxBytes는 전송한 총 바이트를 나타냅니다
	NetworkTxBytes int64 `json:"network_tx_bytes"`
	// BlockRead는 읽은 블록 디바이스 데이터량을 나타냅니다
	BlockRead int64 `json:"block_read"`
	// BlockWrite는 쓴 블록 디바이스 데이터량을 나타냅니다
	BlockWrite int64 `json:"block_write"`
	// PIDs는 컨테이너 내 프로세스 수를 나타냅니다
	PIDs int `json:"pids"`
	// Restarts는 재시작 횟수를 나타냅니다
	Restarts int `json:"restarts"`
	// Labels는 컨테이너 레이블을 나타냅니다
	Labels []ContainerLabel `json:"labels,omitempty"`
	// Health는 컨테이너 상태 정보를 나타냅니다
	Health ContainerHealth `json:"container_health"`
	// Ports는 포트 매핑 정보를 나타냅니다
	Ports []ContainerPort `json:"container_ports"`
	// Networks는 네트워크 설정 정보를 나타냅니다
	Networks []ContainerNetwork `json:"container_network"`
	// Volumes는 볼륨 마운트 정보를 나타냅니다
	Volumes []ContainerVolume `json:"container_volumes"`
	// Env는 환경 변수 정보를 나타냅니다
	Env []ContainerEnv `json:"container_env"`
}

// ContainerLabel은 컨테이너 레이블 정보를 포함하는 구조체입니다.
type ContainerLabel struct {
	Key   string `json:"label_key"`
	Value string `json:"label_value"`
}

// ContainerHealth는 컨테이너 상태 정보를 포함하는 구조체입니다.
type ContainerHealth struct {
	Status          string `json:"status"`
	FailingStreak   int    `json:"failing_streak"`
	LastCheckOutput string `json:"last_check_output"`
}

// ContainerPort는 컨테이너 포트 매핑 정보를 포함하는 구조체입니다.
type ContainerPort struct {
	ContainerPort string `json:"container_port"`
	HostPort      string `json:"host_port"`
	Protocol      string `json:"protocol"`
}

// ContainerNetwork는 컨테이너 네트워크 설정 정보를 포함하는 구조체입니다.
type ContainerNetwork struct {
	Name    string `json:"network_name"`
	IP      string `json:"network_ip"`
	MAC     string `json:"network_mac"`
	RxBytes string `json:"network_rx_bytes"`
	TxBytes string `json:"network_tx_bytes"`
}

// ContainerVolume은 컨테이너 볼륨 마운트 정보를 포함하는 구조체입니다.
type ContainerVolume struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	Type        string `json:"type"`
}

// ContainerEnv는 컨테이너 환경 변수 정보를 포함하는 구조체입니다.
type ContainerEnv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ServiceInfo는 시스템 서비스의 상태 정보를 포함하는 구조체입니다
type ServiceInfo struct {
	// Name은 서비스 이름을 나타냅니다
	Name string `json:"name"`
	// Status는 서비스의 현재 상태를 나타냅니다
	Status string `json:"status"`
	// IsRunning은 서비스 실행 여부를 나타냅니다
	IsRunning bool `json:"is_running"`
}

// ToJSON은 주어진 인터페이스를 JSON 바이트 배열로 변환합니다.
func ToJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling to JSON: %v", err)
		return []byte("{}")
	}
	return data
}
