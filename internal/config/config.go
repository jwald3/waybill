package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Database struct {
		Host            string
		Port            string
		User            string
		Password        string
		DBName          string
		SSLMode         string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}

	App struct {
		Environment string
		LogLevel    string
		APIKey      string
		APIVersion  string
		CORSOrigins []string
		DebugMode   bool
	}

	RateLimit struct {
		Enabled  bool
		Requests int
		Duration time.Duration
	}
}

func Load() *Config {
	config := &Config{}

	config.Server.Port = getEnv("SERVER_PORT", "8000")
	config.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	config.Server.ReadTimeout = getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second)
	config.Server.WriteTimeout = getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second)
	config.Server.IdleTimeout = getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second)

	config.Database.Host = getEnv("DB_HOST", "")
	config.Database.Port = getEnv("DB_PORT", "")
	config.Database.User = getEnv("DB_USER", "admin")
	config.Database.Password = getEnv("DB_PASSWORD", "")
	config.Database.DBName = getEnv("DB_NAME", "myapp")

	config.App.Environment = getEnv("APP_ENV", "development")
	config.App.APIKey = getEnv("API_KEY", "ABC123")
	config.App.LogLevel = getEnv("LOG_LEVEL", "info")
	config.App.APIVersion = getEnv("API_VERSION", "v1")
	config.App.DebugMode = getBoolEnv("DEBUG_MODE", true)
	config.App.CORSOrigins = getSliceEnv("CORS_ORIGINS", []string{"http://localhost:3000"})

	config.RateLimit.Enabled = getBoolEnv("RATE_LIMIT_ENABLED", true)
	config.RateLimit.Requests = getIntEnv("RATE_LIMIT_REQUESTS", 100)
	config.RateLimit.Duration = getDurationEnv("RATE_LIMIT_DURATION", time.Hour)

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=Waybill",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
	)
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
