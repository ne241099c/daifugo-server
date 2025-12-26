package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	JWTSecret     string
	AllowedOrigin string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
}

// Load は環境変数から設定を読み込む
func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-me"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "*"),

		DBUser:     getEnv("DB_USER", "daifugo"),
		DBPassword: getEnv("DB_PASSWORD", "daifugo_pass"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "daifugo_db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}
