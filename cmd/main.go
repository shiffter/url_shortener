package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url_shortener/config"
	"url_shortener/internal/server"
)

func main() {
	v, err := config.LoadConfig()

	logger := logrus.New()
	if err != nil {
		log.Fatal("Cannot load config: ", err.Error())
	}
	cfg, err := config.ParseConfig(v)

	srv := server.NewServer(cfg, logger)
	go srv.MustRunHttp()
	go srv.MustRunGrpc()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	srv.Stop()
	logger.Infof("app stop")
}
