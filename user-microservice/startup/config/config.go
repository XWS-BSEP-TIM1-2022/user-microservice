package config

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port            string
	UserDBHost      string
	UserDBPort      string
	UserServiceName string
	ExpiresIn       time.Duration
	CommonPasswords []string
}

func NewConfig() *Config {
	return &Config{
		Port:            getEnv("USER_SERVICE_PORT", "8085"),
		UserDBHost:      getEnv("USER_DB_HOST", "dislinkt:WiYf6BvFmSpJS2Ob@xws.cjx50.mongodb.net/usersDB"),
		UserDBPort:      getEnv("USER_DB_PORT", ""),
		UserServiceName: getEnv("USER_SERVICE_NAME", "user_service"),
		ExpiresIn:       30 * time.Minute,
		CommonPasswords: getPasswords(),
	}
}

func getPasswords() []string {

	file, err := os.Open("common_passwords.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var retVal []string

	for scanner.Scan() {
		retVal = append(retVal, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	return retVal
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
