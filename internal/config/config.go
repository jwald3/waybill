/// config.go is a file that provides the primary, persistent settings for the application. We load the configurations in cmd/api/main.go primarily to ensure that services can easily import settings.
/// This file loads in environment variables with default fallback values. Feel free to change defaults, but ensure that you do not expose sensitive values in version controlled files.
/// If you save the env variables securely and match the key names, you can leverage this file safely.

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// include any additional properties needed for utilizing the server itself
	Server struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	// these database configurations should be mostly standard across SQL dialects, but any
	// additional information needed for establishing an effective connection with the db should be stored here
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

	// application-specific logic, i.e., what version of the app you're running on
	// this is less about the database or server and more about the actual application itself
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

	// feel free to add any additional service settings, such as Redis, message queueing, etc.
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

// string `os.LookupEnv` wrapper with support for default value specs
func getEnv(key, defaultValue string) string {
	// use `LookupEnv` to attempt to find the env variable. The method returns the value (or nil) if it can find it and returns a bool based on finding the requested key
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// time.Duration `os.LookupEnv` wrapper with support for default value specs
/*
	The env variable could be something like `2h` or `30m` or `1000ms` or something like that
	this method will take that value, convert it into a `time` friendly format

	here are all the time units:
	* "ns"
	* "us"
	* "µs"
	* "μs"
	* "ms"
	* "s"
	* "m"
	* "h"
*/
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// int `os.LookupEnv` wrapper with support for default value specs
/*
	It would make sense to format it like export RATE_LIMIT_REQUESTS=100 where this would easily ensure that it isn't
	trying to work with "100" as a string or something.
*/
func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// bool `os.LookupEnv` wrapper with support for default value specs
/*
	This method works with t, true, 1, f, false, 0 (case insensitive)
	you can do something like export DEBUG_MODE=true and that will work as expected
	(alternatively, you could do DEBUG_MODE=F and that would return false as the value)
*/
func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// slice `os.LookupEnv` wrapper with support for default value specs
/*
	This method takes in a comma-separated string and turns it into a slice of individual strings
	export CORS_ORIGINS="http://localhost:3000,http://website.com" would be a slice of len 2

*/
func getSliceEnv(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// creates the database connection string using the config object
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
