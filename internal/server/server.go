package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/stdlib" // pgx driver
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"url_shortener/config"
)

type Server struct {
	log   *log.Logger
	cfg   *config.Config
	fiber *fiber.App
	grpc  *grpc.Server
}

func NewServer(cfg *config.Config, log *log.Logger) *Server {

	return &Server{
		fiber: fiber.New(fiber.Config{
			AppName:               "UrlShortener",
			DisableStartupMessage: true,
		}),
		cfg:  cfg,
		log:  log,
		grpc: grpc.NewServer(),
	}
}

func (s *Server) MustRunHttp() {
	if err := s.RunHttp(); err != nil {
		panic(err)
	}
}

func (s *Server) MustRunGrpc() {
	if err := s.RunGrpc(); err != nil {
		panic(err)
	}
}

func (s *Server) RunHttp() error {
	if err := s.BuildSrv(s.fiber, s.log); err != nil {
		s.log.Fatalf("Cannot map delivery: %s", err.Error())
	}
	s.log.Infof("Start server on port: %s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	if err := s.fiber.Listen(fmt.Sprintf(":%s", s.cfg.Server.Port)); err != nil {
		s.log.Fatalf("Error starting Server: %s", err.Error())
	}
	return nil
}

func (s *Server) RunGrpc() error {
	l, err := net.Listen("tcp", "localhost:6969")
	if err != nil {
		s.log.Fatalf("Error net listen grpc: %s", err.Error())
	}

	s.log.Infof("grpc start %s", l.Addr().String())
	if err = s.grpc.Serve(l); err != nil {
		s.log.Fatalf("Error serve grpc: %s", err.Error())
	}
	return nil
}

func (s *Server) Stop() {
	err := s.fiber.Shutdown()
	if err != nil {
		s.log.Error(err.Error())
	}
}
