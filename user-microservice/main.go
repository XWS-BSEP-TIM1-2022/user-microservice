package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"user-microservice/application"
	"user-microservice/startup"
	cfg "user-microservice/startup/config"
)

var log = logrus.New()

func main() {
	application.Log = log
	log.Out = os.Stdout

	path := "user-microservice.log"

	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(8760)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	if err == nil {
		log.SetOutput(writer)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.Info("Server starting...")
	
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	defer server.Stop()
}
