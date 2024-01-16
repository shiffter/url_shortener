package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"url_shortener/internal/delivery/http"
	"url_shortener/internal/storage"
	"url_shortener/internal/storage/in_memory"
	"url_shortener/internal/storage/postgres"

	"url_shortener/internal/usecase/links_usecase"
)

func (s *Server) BuildSrv(app *fiber.App, log *logrus.Logger) error {

	storageData := storage.Storage(nil)
	if s.cfg.Server.StorageMode == "postgres" {
		storageData = postgres.NewPostgresStorage(log, s.cfg)
	} else {
		storageData = in_memory.NewMemoryStorage(log)
	}

	linksUsecases := links_usecase.NewLinksUsecase(log, storageData)
	linksHandlers := http.NewLinksHandler(log, linksUsecases)
	http.MapLinksRoutes(app, linksHandlers)
	return nil
}
