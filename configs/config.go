package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	InfluxDB struct {
		URL    string `yaml:"url"`
		Token  string `yaml:"token"`
		Org    string `yaml:"org"`
		Bucket string `yaml:"bucket"`
	} `yaml:"influxdb"`
	Collector struct {
		Interval int `yaml:"interval"`
	} `yaml:"collector"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"postgres"`
}

var cfg *Config

// Load 함수는 설정 파일을 로드하고 전역 설정 객체를 초기화합니다
func Load(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	cfg = &Config{}
	err = yaml.Unmarshal(buf, cfg)
	if err != nil {
		return err
	}

	// 환경 변수로 설정 오버라이드
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Postgres.Host = dbHost
	}

	if influxURL := os.Getenv("INFLUXDB_URL"); influxURL != "" {
		cfg.InfluxDB.URL = influxURL
	}

	return nil
}

// Get 함수는 현재 로드된 설정을 반환합니다
func Get() *Config {
	return cfg
}
