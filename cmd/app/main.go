package main

import (
	"fmt"
	"log"

	"github.com/restlesswhy/video-merger/config"
	"github.com/restlesswhy/video-merger/internal/server"
	"github.com/restlesswhy/video-merger/pkg/logger"
)

func main() {
	log.Println("Starting microservice")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Named(fmt.Sprintf(`(%s)`, cfg.ServiceName))
	appLogger.Infof("CFG: %+v", cfg)

	if err := server.New(appLogger, cfg, nil).Run(); err != nil {
		appLogger.Fatal(err)
	}
}
