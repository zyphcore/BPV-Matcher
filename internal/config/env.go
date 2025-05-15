package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	Environment    string
	JWTSecret      string
	JWTExpiration  time.Duration
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	AllowedOrigins string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	config := Config{
		ServerPort:     "8080",
		Environment:    "development",
		JWTSecret:      "your-secret-key-change-in-production",
		JWTExpiration:  24 * time.Hour,
		DBHost:         "localhost",
		DBPort:         "5432",
		DBUser:         "postgres",
		DBPassword:     "postgres",
		DBName:         "bpvmatcher",
		AllowedOrigins: "*",
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.ServerPort = port
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWTSecret = secret
	}

	if exp := os.Getenv("JWT_EXPIRATION_HOURS"); exp != "" {
		hours, err := strconv.Atoi(exp)
		if err == nil && hours > 0 {
			config.JWTExpiration = time.Duration(hours) * time.Hour
		} else {
			log.Printf("Invalid JWT_EXPIRATION_HOURS: %s, using default", exp)
		}
	}

	if host := os.Getenv("DB_HOST"); host != "" {
		config.DBHost = host
	}

	if port := os.Getenv("DB_PORT"); port != "" {
		config.DBPort = port
	}

	if user := os.Getenv("DB_USER"); user != "" {
		config.DBUser = user
	}

	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.DBPassword = password
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.DBName = dbName
	}

	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		config.AllowedOrigins = origins
	}

	return config
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
