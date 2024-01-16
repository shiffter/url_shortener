package http

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"url_shortener/internal/customError"
	"url_shortener/internal/delivery"
	"url_shortener/internal/usecase"
)

type LinksHandler struct {
	log     *logrus.Logger
	linksUC usecase.LinksUsecase
}

func NewLinksHandler(log *logrus.Logger, linksUC usecase.LinksUsecase) *LinksHandler {
	return &LinksHandler{
		log:     log,
		linksUC: linksUC,
	}
}

func (lh *LinksHandler) CreateShortUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		createShortUrlRequest := delivery.CreateShortUrlRequest{}
		bytesBody := ctx.Body()
		err := json.Unmarshal(bytesBody, &createShortUrlRequest)
		if err != nil {
			lh.log.Error(err)
			sendStatusErr := ctx.SendStatus(http.StatusBadRequest)
			if sendStatusErr != nil {
				lh.log.Error(sendStatusErr.Error())
			}
			return ctx.JSON(delivery.CreateShortUrlResponse{Status: http.StatusBadRequest, Error: err.Error()})
		}
		if len(createShortUrlRequest.OriginalUrl) == 0 {
			err := customError.ErrBadRequest
			lh.log.Error(err)
			sendStatusErr := ctx.SendStatus(http.StatusBadRequest)
			if sendStatusErr != nil {
				lh.log.Error(sendStatusErr.Error())
			}
			return ctx.JSON(delivery.CreateShortUrlResponse{Status: http.StatusBadRequest, Error: err.Error()})
		}

		shortUrl, err := lh.linksUC.CreateShortUrl(context.Background(), createShortUrlRequest.OriginalUrl)
		if err != nil {
			lh.log.Error(err)
			if err.Error() == customError.UrlAlreadyExist.Error() {
				sendStatusErr := ctx.SendStatus(http.StatusConflict)
				if sendStatusErr != nil {
					lh.log.Error(sendStatusErr.Error())
				}
				return ctx.JSON(delivery.CreateShortUrlResponse{Status: http.StatusConflict, Error: err.Error()})
			}

			return ctx.JSON(delivery.CreateShortUrlResponse{Status: http.StatusInternalServerError,
				Error: customError.InternalErr.Error()})
		}

		return ctx.JSON(delivery.CreateShortUrlResponse{ShortUrl: shortUrl, Status: http.StatusOK})
	}
}

func (lh *LinksHandler) GetOriginalUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		shortUrl := ctx.Query("short_url")
		if len(shortUrl) != 10 {
			err := customError.ErrBadRequest
			lh.log.Error(err)
			sendStatusErr := ctx.SendStatus(http.StatusBadRequest)
			if sendStatusErr != nil {
				lh.log.Error(sendStatusErr.Error())
			}
			return ctx.JSON(delivery.GetOriginalUrlResponse{Status: http.StatusBadRequest, Error: err.Error()})
		}

		originalUrl, err := lh.linksUC.GetOriginalUrl(context.Background(), shortUrl)
		if err != nil {
			lh.log.Error(err)
			if err.Error() == "no rows in result set" {
				sendStatusErr := ctx.SendStatus(http.StatusNotFound)
				if sendStatusErr != nil {
					lh.log.Error(sendStatusErr.Error())
				}
				return ctx.JSON(delivery.GetOriginalUrlResponse{Status: http.StatusNotFound,
					Error: customError.NotFound.Error()})
			} else {
				return ctx.JSON(delivery.GetOriginalUrlResponse{Status: http.StatusInternalServerError,
					Error: customError.InternalErr.Error()})
			}
		}

		return ctx.JSON(delivery.GetOriginalUrlResponse{OriginalUrl: originalUrl, Status: http.StatusOK})
	}
}
