package storage

import (
	"context"
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
		"key": metrics.Key, // SystemMetrics 구조체에 Key 필드 추가 필요
	}

	// CPU 메트릭 저장
	cpuFields := map[string]interface{}{
		"usage_percent": metrics.CPU.UsagePercent,
		"temperature":   metrics.CPU.Temperature,
		"load_avg_1":    metrics.CPU.LoadAvg1,
		"load_avg_5":    metrics.CPU.LoadAvg5,
		"load_avg_15":   metrics.CPU.LoadAvg15,
	}
	// 코어별 CPU 사용률 저장
	for core, usage := range metrics.CPU.PerCorePercent {
		cpuFields["core_"+core] = usage
	}

	if err := i.WriteMetrics("cpu", tags, cpuFields); err != nil {
		return err
	}

	// 메모리 메트릭 저장
	memoryFields := map[string]interface{}{
		"total":         metrics.Memory.Total,
		"used":          metrics.Memory.Used,
		"free":          metrics.Memory.Free,
		"usage_percent": metrics.Memory.UsagePercent,
	}
	if err := i.WriteMetrics("memory", tags, memoryFields); err != nil {
		return err
	}

	// 디스크 메트릭 저장
	for _, disk := range metrics.Disk {
		diskTags := map[string]string{
			"key":      metrics.Key,
			"device":   disk.Device,
			"mount_pt": disk.MountPoint,
		}
		diskFields := map[string]interface{}{
			"total":         disk.Total,
			"used":          disk.Used,
			"free":          disk.Free,
			"usage_percent": disk.UsagePercent,
		}

		if err := i.WriteMetrics("disk", diskTags, diskFields); err != nil {
			return err
		}
	}

	// 네트워크 메트릭 저장
	for _, network := range metrics.Network {
		networkTags := map[string]string{
			"key":       metrics.Key,
			"interface": network.Interface,
		}
		networkFields := map[string]interface{}{
			"rx_bytes":   network.RxBytes,
			"tx_bytes":   network.TxBytes,
			"rx_packets": network.RxPackets,
			"tx_packets": network.TxPackets,
			"rx_errors":  network.RxErrors,
			"tx_errors":  network.TxErrors,
			"rx_dropped": network.RxDropped,
			"tx_dropped": network.TxDropped,
		}

		if err := i.WriteMetrics("network", networkTags, networkFields); err != nil {
			return err
		}
	}

	// 프로세스 메트릭 저장
	for _, process := range metrics.Processes {
		processTags := map[string]string{
			"key":  metrics.Key,
			"pid":  strconv.Itoa(process.PID),
			"name": process.Name,
		}

		processFields := map[string]interface{}{
			"cpu":          process.CPU,
			"memory":       process.Memory,
			"user":         process.User,
			"command":      process.Command,
			"thread_count": process.ThreadCount,
			"status":       process.Status,
		}

		if err := i.WriteMetrics("process", processTags, processFields); err != nil {
			return err
		}
	}

	// Docker 컨테이너 메트릭 저장
	for _, container := range metrics.Containers {
		containerTags := map[string]string{
			"key":            metrics.Key,
			"container_id":   container.ID,
			"container_name": container.Name,
			"image":          container.Image,
		}
		containerFields := map[string]interface{}{
			"cpu":         container.CPU,
			"memory":      container.Memory,
			"network_rx":  container.NetworkRx,
			"network_tx":  container.NetworkTx,
			"block_read":  container.BlockRead,
			"block_write": container.BlockWrite,
			"status":      container.Status,
		}
		if err := i.WriteMetrics("docker", containerTags, containerFields); err != nil {
			return err
		}
	}

	// 모든 데이터 저장 후 단일 WriteAPI 인스턴스로 Flush
	i.writeAPI.Flush()

	elapsed := time.Since(start)
	log.Printf("InfluxDB 저장 완료: %v ms 소요 (메트릭 키: %s)", elapsed.Milliseconds(), metrics.Key)

	return nil
}
