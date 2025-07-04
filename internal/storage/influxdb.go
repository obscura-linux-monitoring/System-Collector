package storage

import (
	"context"
	"fmt"

	config "system-collector/configs"
	"system-collector/pkg/logger"
	"system-collector/pkg/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxDBClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	org      string
	bucket   string
	// mutex    sync.Mutex
}

func NewInfluxDBClient() (*InfluxDBClient, error) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("InfluxDBClient 초기화 중")

	cfg := config.Get()

	options := influxdb2.DefaultOptions()
	options.SetBatchSize(500)
	options.SetFlushInterval(1000)
	options.SetUseGZip(true)
	options.SetHTTPRequestTimeout(5000)
	options.SetRetryInterval(1000)
	options.SetMaxRetries(3)
	options.SetLogLevel(0) // 로깅 레벨 낮춤

	client := influxdb2.NewClientWithOptions(cfg.InfluxDB.URL, cfg.InfluxDB.Token, options)

	// 연결 테스트
	_, err := client.Ping(context.Background())
	if err != nil {
		sugar.Errorw("InfluxDB 연결 테스트 실패", "error", err)
		return nil, err
	}

	// 비동기 쓰기 API 사용
	writeAPI := client.WriteAPI(cfg.InfluxDB.Org, cfg.InfluxDB.Bucket)

	// 에러 처리
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			sugar.Errorw("InfluxDB 쓰기 오류", "error", err)
		}
	}()

	return &InfluxDBClient{
		client:   client,
		writeAPI: writeAPI,
		org:      cfg.InfluxDB.Org,
		bucket:   cfg.InfluxDB.Bucket,
	}, nil
}

func (i *InfluxDBClient) WritePoints(points []*write.Point) {
	sugar := logger.GetCustomLogger()
	sugar.Infow("InfluxDBClient 쓰기 시작")

	for _, p := range points {
		i.writeAPI.WritePoint(p)
	}

	sugar.Infow("InfluxDBClient 쓰기 완료")
}

func (i *InfluxDBClient) Close() {
	sugar := logger.GetCustomLogger()
	sugar.Infow("InfluxDBClient 종료 중")

	// 닫기 전에 남은 데이터 플러시
	i.writeAPI.Flush()
	i.client.Close()

	sugar.Infow("InfluxDBClient 종료 완료")
}

func (i *InfluxDBClient) StoreMetrics(metrics *models.SystemMetrics) error {
	sugar := logger.GetCustomLogger()
	sugar.Infow("InfluxDBClient 메트릭스 저장 시작")

	points := make([]*write.Point, 0, 100) // 예상 포인트 수로 초기화

	// 시스템 기본 태그
	baseTags := map[string]string{
		"node_id":     metrics.Key,
		"obscura_key": metrics.USER_ID,
		"hostname":    metrics.System.Hostname,
	}

	// 1. CPU 메트릭스
	cpuFields := map[string]interface{}{
		"architecture":        metrics.CPU.Architecture,
		"model":               metrics.CPU.Model,
		"vendor":              metrics.CPU.Vendor,
		"cache_size":          metrics.CPU.CacheSize,
		"clock_speed":         metrics.CPU.ClockSpeed,
		"total_cores":         metrics.CPU.TotalCores,
		"total_logical_cores": metrics.CPU.TotalLogicalCores,
		"usage":               metrics.CPU.Usage,
		"temperature":         metrics.CPU.Temperature,
		"has_vmx":             metrics.CPU.HasVMX,
		"has_svm":             metrics.CPU.HasSVM,
		"has_avx":             metrics.CPU.HasAVX,
		"has_avx2":            metrics.CPU.HasAVX2,
		"has_neon":            metrics.CPU.HasNEON,
		"has_sve":             metrics.CPU.HasSVE,
		"is_hyperthreading":   metrics.CPU.IsHyperthreading,
		"l1_cache_size":       metrics.CPU.L1CacheSize,
		"l2_cache_size":       metrics.CPU.L2CacheSize,
		"l3_cache_size":       metrics.CPU.L3CacheSize,
		"base_clock_speed":    metrics.CPU.BaseClockSpeed,
		"max_clock_speed":     metrics.CPU.MaxClockSpeed,
		"min_clock_speed":     metrics.CPU.MinClockSpeed,
	}

	// CPU 코어별 정보 추가
	for _, core := range metrics.CPU.Cores {
		cpuFields[fmt.Sprintf("core_%d_usage", core.ID)] = core.Usage
		cpuFields[fmt.Sprintf("core_%d_temperature", core.ID)] = core.Temperature
	}

	cpuPoint := influxdb2.NewPoint("cpu", baseTags, cpuFields, metrics.Timestamp)
	points = append(points, cpuPoint)

	// 2. 메모리 메트릭스
	memFields := map[string]interface{}{
		"total":               metrics.Memory.Total,
		"used":                metrics.Memory.Used,
		"free":                metrics.Memory.Free,
		"available":           metrics.Memory.Available,
		"buffers":             metrics.Memory.Buffers,
		"cached":              metrics.Memory.Cached,
		"swap_total":          metrics.Memory.SwapTotal,
		"swap_used":           metrics.Memory.SwapUsed,
		"swap_free":           metrics.Memory.SwapFree,
		"usage_percent":       metrics.Memory.UsagePercent,
		"data_rate":           metrics.Memory.DataRate,
		"using_slot_count":    metrics.Memory.UsingSlotCount,
		"total_slot_count":    metrics.Memory.TotalSlotCount,
		"form_factor":         metrics.Memory.FormFactor,
		"paged_pool_size":     metrics.Memory.PagedPoolSize,
		"non_paged_pool_size": metrics.Memory.NonPagedPoolSize,
	}
	memPoint := influxdb2.NewPoint("memory", baseTags, memFields, metrics.Timestamp)
	points = append(points, memPoint)

	// 3. 디스크 메트릭스
	for _, disk := range metrics.Disk {
		diskTags := map[string]string{
			"node_id":         metrics.Key,
			"obscura_key":     metrics.USER_ID,
			"device":          disk.Device,
			"mount_point":     disk.MountPoint,
			"filesystem_type": disk.FilesystemType,
		}

		diskFields := map[string]interface{}{
			"total":             disk.Total,
			"used":              disk.Used,
			"free":              disk.Free,
			"usage_percent":     disk.UsagePercent,
			"inodes_total":      disk.InodesTotal,
			"inodes_used":       disk.InodesUsed,
			"inodes_free":       disk.InodesFree,
			"error_flag":        disk.ErrorFlag,
			"error_message":     disk.ErrorMessage,
			"model_name":        disk.ModelName,
			"type":              disk.Type,
			"is_system_disk":    disk.IsSystemDisk,
			"is_page_file_disk": disk.IsPageFileDisk,
		}

		// IO 통계 추가
		diskFields["read_bytes"] = disk.IOStats.ReadBytes
		diskFields["write_bytes"] = disk.IOStats.WriteBytes
		diskFields["reads"] = disk.IOStats.Reads
		diskFields["writes"] = disk.IOStats.Writes
		diskFields["read_bytes_per_sec"] = disk.IOStats.ReadBytesPerSec
		diskFields["write_bytes_per_sec"] = disk.IOStats.WriteBytesPerSec
		diskFields["reads_per_sec"] = disk.IOStats.ReadsPerSec
		diskFields["writes_per_sec"] = disk.IOStats.WritesPerSec
		diskFields["io_in_progress"] = disk.IOStats.IOInProgress
		diskFields["io_time"] = disk.IOStats.IOTime
		diskFields["read_time"] = disk.IOStats.ReadTime
		diskFields["write_time"] = disk.IOStats.WriteTime
		diskFields["error_flag"] = disk.IOStats.ErrorFlag

		diskPoint := influxdb2.NewPoint("disk", diskTags, diskFields, metrics.Timestamp)
		points = append(points, diskPoint)
	}

	// 4. 네트워크 메트릭스
	for _, network := range metrics.Network {
		netTags := map[string]string{
			"node_id":     metrics.Key,
			"obscura_key": metrics.USER_ID,
			"interface":   network.Interface,
			"mac":         network.MAC,
		}

		netFields := map[string]interface{}{
			"ipv4":            network.IPv4,
			"ipv6":            network.IPv6,
			"mtu":             network.MTU,
			"speed":           network.Speed,
			"status":          network.Status,
			"rx_bytes":        network.RxBytes,
			"tx_bytes":        network.TxBytes,
			"rx_packets":      network.RxPackets,
			"tx_packets":      network.TxPackets,
			"rx_errors":       network.RxErrors,
			"tx_errors":       network.TxErrors,
			"rx_dropped":      network.RxDropped,
			"tx_dropped":      network.TxDropped,
			"rx_bytes_sec":    network.RxBytesPerSec,
			"tx_bytes_sec":    network.TxBytesPerSec,
			"connection_type": network.ConnectionType,
		}

		netPoint := influxdb2.NewPoint("network", netTags, netFields, metrics.Timestamp)
		points = append(points, netPoint)
	}

	// 5. 프로세스 메트릭스
	for _, process := range metrics.Processes {

		procTags := map[string]string{
			"node_id":     metrics.Key,
			"obscura_key": metrics.USER_ID,
			"pid":         fmt.Sprintf("%d", process.PID),
			"name":        process.Name,
			"user":        process.User,
			"command":     process.Command,
		}

		procFields := map[string]interface{}{
			"ppid":           process.PPID,
			"status":         process.Status,
			"cpu_time":       process.CPUTime,
			"cpu_usage":      process.CPUUsage,
			"memory_rss":     process.MemoryRSS,
			"memory_vsz":     process.MemoryVSZ,
			"nice":           process.Nice,
			"threads":        process.Threads,
			"open_files":     process.OpenFiles,
			"start_time":     process.StartTime,
			"io_read_bytes":  process.IOReadBytes,
			"io_write_bytes": process.IOWriteBytes,
		}

		procPoint := influxdb2.NewPoint("process", procTags, procFields, metrics.Timestamp)
		points = append(points, procPoint)
	}

	// 6. 컨테이너 메트릭스
	for _, container := range metrics.Containers {
		containerTags := map[string]string{
			"node_id":     metrics.Key,
			"obscura_key": metrics.USER_ID,
			"id":          container.ID,
			"name":        container.Name,
			"image":       container.Image,
		}

		containerFields := map[string]interface{}{
			"status":           container.Status,
			"created":          container.Created,
			"cpu_usage":        container.CPUUsage,
			"memory_usage":     container.MemoryUsage,
			"memory_limit":     container.MemoryLimit,
			"memory_percent":   container.MemoryPercent,
			"network_rx_bytes": container.NetworkRxBytes,
			"network_tx_bytes": container.NetworkTxBytes,
			"block_read":       container.BlockRead,
			"block_write":      container.BlockWrite,
			"pids":             container.PIDs,
			"restarts":         container.Restarts,
			// Health 정보 추가
			"health_status":            container.Health.Status,
			"health_failing_streak":    container.Health.FailingStreak,
			"health_last_check_output": container.Health.LastCheckOutput,
		}

		// 레이블 정보 추가
		for i, label := range container.Labels {
			containerFields[fmt.Sprintf("label_%d_key", i)] = label.Key
			containerFields[fmt.Sprintf("label_%d_value", i)] = label.Value
		}

		// 포트 정보 추가
		for i, port := range container.Ports {
			containerFields[fmt.Sprintf("port_%d_container", i)] = port.ContainerPort
			containerFields[fmt.Sprintf("port_%d_host", i)] = port.HostPort
			containerFields[fmt.Sprintf("port_%d_protocol", i)] = port.Protocol
		}

		// 네트워크 정보 추가
		for i, network := range container.Networks {
			containerFields[fmt.Sprintf("network_%d_name", i)] = network.Name
			containerFields[fmt.Sprintf("network_%d_ip", i)] = network.IP
			containerFields[fmt.Sprintf("network_%d_mac", i)] = network.MAC
			containerFields[fmt.Sprintf("network_%d_rx_bytes", i)] = network.RxBytes
			containerFields[fmt.Sprintf("network_%d_tx_bytes", i)] = network.TxBytes
		}

		// 볼륨 정보 추가
		for i, volume := range container.Volumes {
			containerFields[fmt.Sprintf("volume_%d_source", i)] = volume.Source
			containerFields[fmt.Sprintf("volume_%d_destination", i)] = volume.Destination
			containerFields[fmt.Sprintf("volume_%d_mode", i)] = volume.Mode
			containerFields[fmt.Sprintf("volume_%d_type", i)] = volume.Type
		}

		// 환경 변수 정보 추가
		for i, env := range container.Env {
			containerFields[fmt.Sprintf("env_%d_key", i)] = env.Key
			containerFields[fmt.Sprintf("env_%d_value", i)] = env.Value
		}

		containerPoint := influxdb2.NewPoint("container", containerTags, containerFields, metrics.Timestamp)
		points = append(points, containerPoint)
	}

	// 7. 서비스 메트릭스
	for _, service := range metrics.Services {
		serviceTags := map[string]string{
			"node_id":     metrics.Key,
			"obscura_key": metrics.USER_ID,
			"name":        service.Name,
		}

		serviceFields := map[string]interface{}{
			"status":       service.Status,
			"enabled":      service.Enabled,
			"type":         service.Type,
			"load_state":   service.LoadState,
			"active_state": service.ActiveState,
			"sub_state":    service.SubState,
		}

		servicePoint := influxdb2.NewPoint("service", serviceTags, serviceFields, metrics.Timestamp)
		points = append(points, servicePoint)
	}

	// 8. 시스템 정보 메트릭스
	systemTags := map[string]string{
		"node_id":     metrics.Key,
		"obscura_key": metrics.USER_ID,
	}

	systemFields := map[string]interface{}{
		"os_name":                metrics.System.OSName,
		"os_version":             metrics.System.OSVersion,
		"os_architecture":        metrics.System.OSArchitecture,
		"os_kernel_version":      metrics.System.OSKernelVersion,
		"boot_time":              metrics.System.BootTime,
		"uptime":                 metrics.System.Uptime,
		"total_processes":        metrics.System.TotalProcesses,
		"total_threads":          metrics.System.TotalThreads,
		"total_file_descriptors": metrics.System.TotalFileDescriptors,
	}

	systemPoint := influxdb2.NewPoint("system", systemTags, systemFields, metrics.Timestamp)
	points = append(points, systemPoint)

	// 모든 포인트를 한 번에 전송
	i.WritePoints(points)

	sugar.Infow("InfluxDBClient 메트릭스 저장 완료")
	return nil
}
