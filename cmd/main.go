package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"url_shortener/config"
	"url_shortener/internal/server"
)

func main() {
	v, err := config.LoadConfig()

	logger := logrus.New()
	if err != nil {
		log.Fatal("Cannot cload config: ", err.Error())
	}
	cfg, err := config.ParseConfig(v)

	fmt.Println(cfg)
	srv := server.NewServer(cfg, logger)
	if err := srv.Run(); err != nil {
		logger.Fatalf("cant start srv %s", err.Error())
	}
}
