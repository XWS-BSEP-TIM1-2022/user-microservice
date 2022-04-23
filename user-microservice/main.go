package main

import (
	"user-microservice/startup"
	cfg "user-microservice/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	defer server.Stop()
}
