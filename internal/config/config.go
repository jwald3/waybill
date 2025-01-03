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

	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnv("DB_PORT", "5432")
	config.Database.User = getEnv("DB_USER", "postgres")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.DBName = getEnv("DB_NAME", "myapp")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	config.Database.MaxOpenConns = getIntEnv("DB_MAX_OPEN_CONNS", 25)
	config.Database.MaxIdleConns = getIntEnv("DB_MAX_IDLE_CONNS", 25)
	config.Database.ConnMaxLifetime = getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	config.App.Environment = getEnv("APP_ENV", "development")
	config.App.LogLevel = getEnv("LOG_LEVEL", "info")
	config.App.APIVersion = getEnv("API_VERSION", "v1")
	config.App.DebugMode = getBoolEnv("DEBUG_MODE", true)
	config.App.CORSOrigins = getSliceEnv("CORS_ORIGINS", []string{"http://localhost:3000"})

	config.RateLimit.Enabled = getBoolEnv("RATE_LIMIT_ENABLED", true)
	config.RateLimit.Requests = getIntEnv("RATE_LIMIT_REQUESTS", 100)
	config.RateLimit.Duration = getDurationEnv("RATE_LIMIT_DURATION", time.Hour)

	return config
}

// string `os.LookupEnv` wrapper with support for default value specs
func getEnv(key, defaultValue string) string {
	// use `LookupEnv` to attempt to find the env variable. The method returns the value (or nil) if it can find it and returns a bool based on finding the requested key
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
