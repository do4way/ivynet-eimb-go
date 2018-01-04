package redisConfig

import (
	"os"
	"strconv"
)

//Config a structure of redis settings.
type Config struct {
	Host     string
	Port     string
	DB       int
	Password string
}

var (
	defaultRedisHost = "127.0.0.1"
	defaultRedisPort = "6379"
)

//Read to read redis configuration from runtime environment.
func Read() *Config {
	redisHost := os.Getenv("REDIS_HOST_NAME")
	if redisHost == "" {
		redisHost = defaultRedisHost
	}
	redisPort := os.Getenv("REDIS_HOST_PORT")
	if redisPort == "" {
		redisPort = defaultRedisPort
	}
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisPassword := os.Getenv("REDIS_PASSWORD")

	return &Config{
		Host:     redisHost,
		Port:     redisPort,
		DB:       redisDB,
		Password: redisPassword,
	}
}
