package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	config "system-collector/configs"
	"system-collector/pkg/models"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	api "github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDBClient struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	org      string
	bucket   string
}

func NewInfluxDBClient() (*InfluxDBClient, error) {
	// 설정에서 로드
	cfg := config.Get()

	// 클라이언트 옵션 설정
	options := influxdb2.DefaultOptions()
	options.SetBatchSize(500)
	options.SetFlushInterval(1000)
	options.SetUseGZip(true)
	options.SetHTTPRequestTimeout(30000)
	options.SetRetryInterval(5000)
	options.SetMaxRetries(5)
	options.SetLogLevel(1)

	client := influxdb2.NewClientWithOptions(cfg.InfluxDB.URL, cfg.InfluxDB.Token, options)

	// 연결 테스트
	_, err := client.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	// 단일 WriteAPI 인스턴스 생성
	writeAPI := client.WriteAPI(cfg.InfluxDB.Org, cfg.InfluxDB.Bucket)

	// 에러 처리
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Printf("InfluxDB write error: %s", err.Error())
		}
	}()

	return &InfluxDBClient{
		client:   client,
		writeAPI: writeAPI,
		org:      cfg.InfluxDB.Org,
		bucket:   cfg.InfluxDB.Bucket,
	}, nil
}

func (i *InfluxDBClient) WriteMetrics(measurement string, tags map[string]string, fields map[string]interface{}) error {
	p := influxdb2.NewPoint(
		measurement,
		tags,
		fields,
		time.Now(),
	)

	// 저장된 단일 WriteAPI 사용
	i.writeAPI.WritePoint(p)

	return nil
}

func (i *InfluxDBClient) Close() {
	i.client.Close()
}

func (i *InfluxDBClient) StoreMetrics(metrics *models.SystemMetrics) error {
	start := time.Now()

	tags := map[string]string{
		"key":      metrics.Key,
		"hostname": metrics.System.Hostname,
	}

	// CPU 메트릭스 저장
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
		"is_hyperthreading":   metrics.CPU.IsHyperthreading,
	}

	// CPU 코어별 정보 저장
	for _, core := range metrics.CPU.Cores {
		cpuFields[fmt.Sprintf("core_%d_usage", core.ID)] = core.Usage
		cpuFields[fmt.Sprintf("core_%d_temperature", core.ID)] = core.Temperature
	}

	if err := i.WriteMetrics("cpu", tags, cpuFields); err != nil {
		return err
	}

	// 메모리 메트릭스 저장
	memoryFields := map[string]interface{}{
		"total":         metrics.Memory.Total,
		"used":          metrics.Memory.Used,
		"free":          metrics.Memory.Free,
		"available":     metrics.Memory.Available,
		"buffers":       metrics.Memory.Buffers,
		"cached":        metrics.Memory.Cached,
		"swap_total":    metrics.Memory.SwapTotal,
		"swap_used":     metrics.Memory.SwapUsed,
		"swap_free":     metrics.Memory.SwapFree,
		"usage_percent": metrics.Memory.UsagePercent,
	}
	if err := i.WriteMetrics("memory", tags, memoryFields); err != nil {
		return err
	}

	// 디스크 메트릭스 저장
	for _, disk := range metrics.Disk {
		diskTags := map[string]string{
			"key":             metrics.Key,
			"device":          disk.Device,
			"mount_point":     disk.MountPoint,
			"filesystem_type": disk.FilesystemType,
		}

		diskFields := map[string]interface{}{
			"total":         disk.Total,
			"used":          disk.Used,
			"free":          disk.Free,
			"inodes_total":  disk.InodesTotal,
			"inodes_used":   disk.InodesUsed,
			"inodes_free":   disk.InodesFree,
			"usage_percent": disk.UsagePercent,
			"error_flag":    disk.ErrorFlag,
			// IO 통계
			"io_read_bytes":     disk.IOStats.ReadBytes,
			"io_write_bytes":    disk.IOStats.WriteBytes,
			"io_reads":          disk.IOStats.Reads,
			"io_writes":         disk.IOStats.Writes,
			"io_reads_per_sec":  disk.IOStats.ReadsPerSec,
			"io_writes_per_sec": disk.IOStats.WritesPerSec,
			"io_in_progress":    disk.IOStats.IOInProgress,
			"io_time":           disk.IOStats.IOTime,
		}

		if err := i.WriteMetrics("disk", diskTags, diskFields); err != nil {
			return err
		}
	}

	// 네트워크 메트릭스 저장
	for _, network := range metrics.Network {
		networkTags := map[string]string{
			"key":       metrics.Key,
			"interface": network.Interface,
			"ip":        network.IP,
			"mac":       network.MAC,
		}

		networkFields := map[string]interface{}{
			"mtu":              network.MTU,
			"speed":            network.Speed,
			"status":           network.Status,
			"rx_bytes":         network.RxBytes,
			"tx_bytes":         network.TxBytes,
			"rx_packets":       network.RxPackets,
			"tx_packets":       network.TxPackets,
			"rx_errors":        network.RxErrors,
			"tx_errors":        network.TxErrors,
			"rx_dropped":       network.RxDropped,
			"tx_dropped":       network.TxDropped,
			"rx_bytes_per_sec": network.RxBytesPerSec,
			"tx_bytes_per_sec": network.TxBytesPerSec,
		}

		if err := i.WriteMetrics("network", networkTags, networkFields); err != nil {
			return err
		}
	}

	// 프로세스 메트릭스 저장
	for _, process := range metrics.Processes {
		processTags := map[string]string{
			"key":     metrics.Key,
			"pid":     strconv.Itoa(process.PID),
			"name":    process.Name,
			"user":    process.User,
			"command": process.Command,
		}

		processFields := map[string]interface{}{
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

		if err := i.WriteMetrics("process", processTags, processFields); err != nil {
			return err
		}
	}

	// 도커 컨테이너 메트릭스 저장
	for _, container := range metrics.Containers {
		containerTags := map[string]string{
			"key":            metrics.Key,
			"container_id":   container.ID,
			"container_name": container.Name,
			"image":          container.Image,
		}

		containerFields := map[string]interface{}{
			"status":         container.Status,
			"created":        container.Created,
			"cpu_usage":      container.CPUUsage,
			"memory_usage":   container.MemoryUsage,
			"memory_limit":   container.MemoryLimit,
			"memory_percent": container.MemoryPercent,
			"network_rx":     container.NetworkRxBytes,
			"network_tx":     container.NetworkTxBytes,
			"block_read":     container.BlockRead,
			"block_write":    container.BlockWrite,
			"pids":           container.PIDs,
			"restarts":       container.Restarts,
		}

		if err := i.WriteMetrics("docker", containerTags, containerFields); err != nil {
			return err
		}
	}

	i.writeAPI.Flush()

	elapsed := time.Since(start)
	log.Printf("InfluxDB 저장 완료: %v ms 소요 (메트릭 키: %s)", elapsed.Milliseconds(), metrics.Key)

	return nil
}
