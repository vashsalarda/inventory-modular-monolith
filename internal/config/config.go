package config

import "os"

type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "3000"),
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "inventory_db"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}