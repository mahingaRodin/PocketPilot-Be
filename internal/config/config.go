package config

import (
	"os"
	// "log"
)

type Config struct {
    DatabaseURL         string
    JWTSecret          string
    Port               string
    RedisURL           string
    AWSAccessKeyID     string
    AWSSecretAccessKey string
    S3Bucket          string
    GoogleVisionAPIKey string
}

func Load() *Config {
    return &Config{
        DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/pocket_pilot_db?sslmode=disable"),
        JWTSecret:          getEnv("JWT_SECRET", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkV4cGVuc2VUcmFja2VyIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"),
        Port:               getEnv("PORT", "8080"),
        RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
        AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
        AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
        S3Bucket:          getEnv("S3_BUCKET", "expense-receipts"),
        GoogleVisionAPIKey: getEnv("GOOGLE_VISION_API_KEY", ""),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}