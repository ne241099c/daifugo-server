package config

import (
	"os"
)

type Config struct {
	Port          string
	JWTSecret     string
	AllowedOrigin string
}

// Load は環境変数から設定を読み込む
func Load() *Config {
	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-me"), // デフォルト値は開発用
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "*"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
