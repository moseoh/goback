package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadConfigFromEnv 환경변수에서 설정을 읽어옵니다
// .env 파일이 있으면 먼저 로드합니다
func LoadConfigFromEnv() *BackupConfig {
	// .env 파일 로드 (있는 경우)
	if err := godotenv.Load(); err != nil {
		// .env 파일이 없어도 계속 진행
		fmt.Println("💡 .env 파일을 찾을 수 없습니다. 환경변수를 사용합니다.")
	}

	config := &BackupConfig{
		Host:        getEnvOrDefault("MYSQL_HOST", "localhost"),
		Port:        getEnvOrDefault("MYSQL_PORT", "3306"),
		Username:    getEnvOrDefault("MYSQL_USERNAME", "root"),
		Password:    getEnvOrDefault("MYSQL_PASSWORD", ""),
		Database:    getEnvOrDefault("MYSQL_DATABASE", ""),
		OutputDir:   getEnvOrDefault("BACKUP_OUTPUT_DIR", "./backups"),
		Workers:     getEnvIntOrDefault("BACKUP_WORKERS", runtime.NumCPU()),
		BatchSize:   getEnvIntOrDefault("BACKUP_BATCH_SIZE", 50000),
		MultiInsert: getEnvIntOrDefault("BACKUP_MULTI_INSERT", 1000),
	}

	// 데이터베이스 이름이 비어있으면 경고
	if config.Database == "" {
		log.Println("⚠️ 데이터베이스 이름이 설정되지 않았습니다. 명령행 인수로 지정해주세요.")
	}

	return config
}

// getEnvOrDefault 환경변수 값을 가져오거나 기본값을 반환합니다
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault 환경변수에서 정수값을 가져오거나 기본값을 반환합니다
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
