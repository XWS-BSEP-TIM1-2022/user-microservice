package config

import "os"

type Config struct {
	Port            string
	UserDBHost      string
	UserDBPort      string
	UserServiceName string
}

func NewConfig() *Config {
	return &Config{
		Port:            getEnv("USER_SERVICE_PORT", "8085"),
		UserDBHost:      getEnv("USER_DB_HOST", "userMicroservice:nuEIm8GkSZbm3MKd@xws.cjx50.mongodb.net/usersDB"),
		UserDBPort:      getEnv("USER_DB_PORT", ""),
		UserServiceName: getEnv("USER_SERVICE_NAME", "user_service"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
