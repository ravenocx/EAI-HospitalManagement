package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
	DBSslMode  string

	DbMaxConn         int
	DbMaxIdleConn     int
	DbMaxLifetimeConn int

	ServerHost        string
	ServerPort        string
	ServerReadTimeout int
}

var configOnce sync.Once
var config Config
var err error

func LoadConfig() (Config, error) {
	configOnce.Do(func() {
		err = godotenv.Load()
		if err != nil {
			err = fmt.Errorf("failed to load env file: %v", err)
			return
		}

		config = Config{
			DbUsername: GetEnv("DB_USERNAME", "postgres"),
			DbPassword: GetEnv("DB_PASSWORD", "postgres"),
			DbHost:     GetEnv("DB_HOST", "127.0.0.1"),
			DbPort:     GetEnv("DB_PORT", "5432"),
			DbName:     GetEnv("DB_NAME", "postgres"),
			DBSslMode:  GetEnv("DB_SSL_MODE", "disable"),
			ServerHost: GetEnv("SERVER_HOST", "127.0.0.1"),
			ServerPort: GetEnv("SERVER_PORT", "3000"),
		}

		config.DbMaxConn, err = strconv.Atoi(GetEnv("DB_MAX_CONNECTIONS", "50"))
		if err != nil {
			err = fmt.Errorf("failed to convert DB_MAX_CONNECTIONS to int: %v", err)
			return
		}

		config.DbMaxIdleConn, err = strconv.Atoi(GetEnv("DB_MAX_IDLE_CONNECTIONS", "30"))
		if err != nil {
			err = fmt.Errorf("failed to convert DB_MAX_IDLE_CONNECTIONS to int: %v", err)
			return
		}

		config.DbMaxLifetimeConn, err = strconv.Atoi(GetEnv("DB_MAX_LIFETIME_CONNECTIONS", "1"))
		if err != nil {
			err = fmt.Errorf("failed to convert DB_MAX_LIFETIME_CONNECTIONS to int: %v", err)
			return
		}

		config.ServerReadTimeout, err = strconv.Atoi(GetEnv("SERVER_READ_TIMEOUT", "60"))
		if err != nil {
			err = fmt.Errorf("failed to convert SERVER_READ_TIMEOUT to int: %v", err)
			return
		}
	})

	return config, err
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
