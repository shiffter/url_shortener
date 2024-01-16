package links_usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
	repoMock "url_shortener/internal/storage/mocks"
)

func TestGetOriginalUrl(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockStorage := repoMock.NewMockStorage(controller)

	ctx := context.Background()
	log := logrus.New()
	Usecase := NewLinksUsecase(log, mockStorage)
	shortUrl := "0kG0-yDn0I"
	originalUrl := "testOrigUrl"

	mockStorage.EXPECT().SaveUrlPair(ctx, shortUrl, originalUrl).Return(nil).Times(1)
	originalUrl, err := Usecase.CreateShortUrl(ctx, originalUrl)
	//require.Equal(t, "", originalUrl)
	require.NoError(t, err)
}

//func TestUrlSevice_CreateShortUrl(t *testing.T) {
//	type mockActions func(r *repoMock.MockStorage, shortLink, originalLink string)
//	ctx := context.Background()
//	inputs := []struct {
//		name        string
//		mockActions mockActions
//		request     *delivery.CreateShortUrlRequest
//		response    *delivery.CreateShortUrlResponse
//		expectedErr error
//	}{
//		{
//			name: "good_1",
//			mockActions: func(r *repoMock.MockStorage, shortLink, originalLink string) {
//				r.EXPECT().SaveUrlPair(ctx, shortLink, originalLink).Return(nil).Times(1)
//			},
//			request:     &delivery.CreateShortUrlRequest{OriginalUrl: "time_ticking"},
//			response:    &delivery.CreateShortUrlResponse{ShortUrl: "UxoXKoz9y3"},
//			expectedErr: nil,
//		},
//	}
//	for _, tt := range inputs {
//		t.Run(tt.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			storage := repoMock.NewMockStorage(ctrl)
//			tt.mockActions(storage, Hasher(tt.request.OriginalUrl), tt.request.OriginalUrl)
//
//			log := logrus.New()
//			linksUC := NewLinksUsecase(log, storage)
//
//			app := fiber.New()
//			linksHandler := http.NewLinksHandler(log, linksUC)
//
//			//linksHandler.CreateShortUrl()
//			app.Post("/create_short_url", linksHandler.CreateShortUrl())
//
//			reqBody, _ := json.Marshal(tt.request)
//			req := httptest.NewRequest("POST", "/create_short_url", bytes.NewBuffer(reqBody))
//			//req.Form = url.Values{"original_url": []string{tt.request.OriginalUrl}}
//			resp, _ := app.Test(req, 30)
//			body, _ := io.ReadAll(resp.Body)
//			respDelivery := delivery.CreateShortUrlResponse{}
//			json.Unmarshal(body, &respDelivery)
//
//			assert.Equal(t, tt.response.ShortUrl, respDelivery.ShortUrl)
//		})
//	}
//}
