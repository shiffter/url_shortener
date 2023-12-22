package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"url_shortener/config"
)

type Server struct {
	log   *log.Logger
	cfg   *config.Config
	fiber *fiber.App
}

func NewServer(cfg *config.Config, log *log.Logger) *Server {
	return &Server{
		fiber: fiber.New(fiber.Config{
			AppName:               "UrlShortener",
			DisableStartupMessage: true,
		}),
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run() error {
	if err := s.BuildSrv(s.fiber, s.log); err != nil {
		s.log.Fatalf("Cannot map delivery: ", err)
	}
	s.log.Infof("Start server on port: %s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	if err := s.fiber.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)); err != nil {
		s.log.Fatalf("Error starting Server: ", err)
	}

	return nil
}
