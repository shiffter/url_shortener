package links_usecase

import (
	"context"
	"testing"
	//"url_shortener/pkg/links_usecase"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CheckExistingOriginalUrl(ctx context.Context, originalUrl string) (bool, error) {
	args := m.Called(ctx, originalUrl)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) GetShortUrl(ctx context.Context, originalUrl string) (string, error) {
	args := m.Called(ctx, originalUrl)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) CheckExistingShortUrl(ctx context.Context, shortUrl string) (bool, error) {
	args := m.Called(ctx, shortUrl)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) SaveUrlPair(ctx context.Context, shortUrl string, originalUrl string) error {
	args := m.Called(ctx, shortUrl, originalUrl)
	return args.Error(0)
}

func (m *MockStorage) GetOriginalUrl(ctx context.Context, shortUrl string) (string, error) {
	args := m.Called(ctx, shortUrl)
	return args.String(0), args.Error(1)
}

func TestCreateShortUrl(t *testing.T) {
	log := logrus.New()
	mockStorage := new(MockStorage)
	usecase := NewLinksUsecase(log, mockStorage)

	mockStorage.On("CheckExistingOriginalUrl", mock.Anything, "http://example.com").Return(false, nil)

	mockStorage.On("GetShortUrl", mock.Anything, "http://example.com").Return("", nil)

	mockStorage.On("CheckExistingShortUrl", mock.Anything, mock.Anything).Return(false, nil)

	mockStorage.On("SaveUrlPair", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	shortUrl, err := usecase.CreateShortUrl(context.Background(), "http://example.com")

	assert.NoError(t, err)
	assert.NotEmpty(t, shortUrl)

}

func TestGetOriginalUrl(t *testing.T) {
	log := logrus.New()
	mockStorage := new(MockStorage)
	usecase := NewLinksUsecase(log, mockStorage)

	mockStorage.On("GetOriginalUrl", mock.Anything, "short123").Return("http://example.com", nil)

	originalUrl, err := usecase.GetOriginalUrl(context.Background(), "short123")

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", originalUrl)

	mockStorage.AssertExpectations(t)
}
