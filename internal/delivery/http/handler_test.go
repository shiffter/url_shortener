package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
	"url_shortener/internal/customError"
	"url_shortener/internal/delivery"
	repoMock "url_shortener/internal/storage/mocks"
	"url_shortener/internal/usecase/links_usecase"
)

func TestUrlSevice_CreateShortUrl(t *testing.T) {
	type mockActions func(r *repoMock.MockStorage, shortLink, originalLink string)
	ctx := context.Background()
	inputs := []struct {
		name        string
		route       string
		mockActions mockActions
		request     *delivery.CreateShortUrlRequest
		response    *delivery.CreateShortUrlResponse
	}{
		{
			name:  "good_1",
			route: "/create_short_url",
			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().SaveUrlPair(ctx, shortLink, originalLink).Return(nil).Times(1)
			},
			request:  &delivery.CreateShortUrlRequest{OriginalUrl: "time_ticking"},
			response: &delivery.CreateShortUrlResponse{Status: 200, Error: "", ShortUrl: "UxoXKoz9y3"},
		},
		{
			name:  "good_2",
			route: "/create_short_url",
			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().SaveUrlPair(ctx, shortLink, originalLink).Return(nil).Times(1)
			},
			request:  &delivery.CreateShortUrlRequest{OriginalUrl: "hook_na_fontan"},
			response: &delivery.CreateShortUrlResponse{Status: 200, Error: "", ShortUrl: "lssTXtIuap"},
		},
		{
			name:        "empty_original_url",
			route:       "/create_short_url",
			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {},
			request:     &delivery.CreateShortUrlRequest{OriginalUrl: ""},
			response:    &delivery.CreateShortUrlResponse{ShortUrl: "", Status: 400, Error: customError.ErrBadRequest.Error()},
		},
		{
			name:  "url_already_exist",
			route: "/create_short_url",
			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().SaveUrlPair(ctx, shortLink, originalLink).Return(customError.UrlAlreadyExist).Times(1)
			},
			request:  &delivery.CreateShortUrlRequest{OriginalUrl: "already_exist"},
			response: &delivery.CreateShortUrlResponse{ShortUrl: "", Status: 409, Error: customError.UrlAlreadyExist.Error()},
		},
		{
			name:  "internal_error",
			route: "/create_short_url",
			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {
				r.EXPECT().SaveUrlPair(ctx, shortLink, originalLink).Return(customError.InternalErr).Times(1)
			},
			request:  &delivery.CreateShortUrlRequest{OriginalUrl: "already_exist"},
			response: &delivery.CreateShortUrlResponse{ShortUrl: "", Status: 500, Error: customError.InternalErr.Error()},
		},
	}
	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := repoMock.NewMockStorage(ctrl)
			tt.mockActions(storage, links_usecase.Hasher(tt.request.OriginalUrl), tt.request.OriginalUrl)

			log := logrus.New()
			linksUC := links_usecase.NewLinksUsecase(log, storage)

			app := fiber.New()
			linksHandler := NewLinksHandler(log, linksUC)

			app.Post(tt.route, linksHandler.CreateShortUrl())

			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", tt.route, bytes.NewBuffer(reqBody))
			resp, _ := app.Test(req, 30)

			body, _ := io.ReadAll(resp.Body)
			respDelivery := delivery.CreateShortUrlResponse{}
			json.Unmarshal(body, &respDelivery)
			assert.Equal(t, tt.response, &respDelivery)
		})
	}
}

func TestUrlService_GetOriginalUrl(t *testing.T) {
	type mockActions func(r *repoMock.MockStorage, shortLink string)
	ctx := context.Background()
	inputs := []struct {
		name          string
		route         string
		mockActions   mockActions
		queryUrlParam string
		response      *delivery.GetOriginalUrlResponse
	}{
		{
			name:  "good_1",
			route: `/get_short_url`,
			mockActions: func(r *repoMock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalUrl(ctx, shortLink).Return("hook_na_fontan", nil).Times(1)
			},
			queryUrlParam: "UxoXKoz9y3",
			response:      &delivery.GetOriginalUrlResponse{Status: 200, Error: "", OriginalUrl: "hook_na_fontan"},
		},
		{
			name:  "good_2",
			route: `/get_short_url`,
			mockActions: func(r *repoMock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalUrl(ctx, shortLink).Return("ozon.best", nil).Times(1)
			},
			queryUrlParam: "7zcVU_r3qN",
			response:      &delivery.GetOriginalUrlResponse{Status: 200, Error: "", OriginalUrl: "ozon.best"},
		},
		{
			name:          "url len < 10",
			route:         `/get_short_url`,
			mockActions:   func(r *repoMock.MockStorage, shortLink string) {},
			queryUrlParam: "7zcVU_r3",
			response:      &delivery.GetOriginalUrlResponse{Status: 400, Error: customError.ErrBadRequest.Error(), OriginalUrl: ""},
		},
		{
			name:          "url len > 10",
			route:         `/get_short_url`,
			mockActions:   func(r *repoMock.MockStorage, shortLink string) {},
			queryUrlParam: "7zcVU_r3ssqwe123asd1",
			response:      &delivery.GetOriginalUrlResponse{Status: 400, Error: customError.ErrBadRequest.Error(), OriginalUrl: ""},
		},
		{
			name:  "no existing url",
			route: `/get_short_url`,
			mockActions: func(r *repoMock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalUrl(ctx, shortLink).Return("", errors.New("no rows in result set")).Times(1)
			},
			queryUrlParam: "7zcVU_r3sw",
			response:      &delivery.GetOriginalUrlResponse{Status: 404, Error: customError.NotFound.Error(), OriginalUrl: ""},
		},
		{
			name:  "internal error",
			route: `/get_short_url`,
			mockActions: func(r *repoMock.MockStorage, shortLink string) {
				r.EXPECT().GetOriginalUrl(ctx, shortLink).Return("", customError.InternalErr).Times(1)
			},
			queryUrlParam: "7zcVU_r3sw",
			response:      &delivery.GetOriginalUrlResponse{Status: 500, Error: customError.InternalErr.Error()},
		},
	}

	for _, tt := range inputs {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := repoMock.NewMockStorage(ctrl)
			tt.mockActions(storage, tt.queryUrlParam)

			log := logrus.New()
			linksUC := links_usecase.NewLinksUsecase(log, storage)
			linksHandler := NewLinksHandler(log, linksUC)

			app := fiber.New()
			app.Get(`/get_short_url`, linksHandler.GetOriginalUrl())

			req := httptest.NewRequest("GET", fmt.Sprintf("%s?short_url=%s", tt.route, tt.queryUrlParam), nil)

			resp, _ := app.Test(req, 30)
			body, _ := io.ReadAll(resp.Body)
			respDelivery := delivery.GetOriginalUrlResponse{}
			json.Unmarshal(body, &respDelivery)
			assert.Equal(t, tt.response, &respDelivery)
		})
	}
}
