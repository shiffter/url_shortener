package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"url_shortener/internal/delivery"
	"url_shortener/internal/usecase"
)

type LinksHandler struct {
	log     *logrus.Logger
	linksUC usecase.LinksUsecase
}

//type LinksHandler interface {
//	Pong(s int) error
//}

func NewLinksHandler(log *logrus.Logger, linksUC usecase.LinksUsecase) *LinksHandler {
	return &LinksHandler{
		log:     log,
		linksUC: linksUC,
	}
}

func (lh *LinksHandler) CreateShortUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		createShortUrlParams := delivery.CreateShortUrlParams{}
		bytesBody := ctx.Body()
		err := json.Unmarshal(bytesBody, &createShortUrlParams)
		if err != nil {
			lh.log.Error(err)
			return err
		}
		shortUrl, err := lh.linksUC.CreateShortUrl(context.Background(), createShortUrlParams.OriginalUrl)
		if err != nil {
			lh.log.Error(err)
			return err
		}
		resp := delivery.CreateShortUrlResponse{ShortUrl: shortUrl}
		return ctx.JSON(resp)
	}
}

func (lh *LinksHandler) GetOriginalUrl() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		shortUrl := ctx.Query("short_url")
		if len(shortUrl) != 10 {
			err := errors.New(fmt.Sprintf("bad request, wrong length=%d short url", len(shortUrl)))
			lh.log.Error(err)
			return err
		}
		originalUrl, err := lh.linksUC.GetOriginalUrl(context.Background(), shortUrl)
		if err != nil {
			lh.log.Error(err)
			if err.Error() == "no rows in result set" {
				return ctx.SendStatus(http.StatusNotFound)
			} else {
				return ctx.SendStatus(http.StatusInternalServerError)
			}
		}
		resp := delivery.GetOriginalUrlResp{OriginalUrl: originalUrl}
		return ctx.JSON(resp)
	}
}
