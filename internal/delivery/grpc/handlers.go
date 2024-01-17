package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	url_shortener "url_shortener/internal/proto_gen/shortener"
)

type Usecase interface {
	CreateShortUrl(ctx context.Context, origUrl string) (string, error)
	GetOriginalUrl(ctx context.Context, shortUrl string) (string, error)
}

type ServerApi struct {
	url_shortener.UnimplementedShortenerServer
	usecase Usecase
}

func Register(grpc *grpc.Server, usecase Usecase) {
	url_shortener.RegisterShortenerServer(grpc, &ServerApi{usecase: usecase})
}

func (s *ServerApi) CreateShortUrl(
	ctx context.Context,
	request *url_shortener.NewShortUrlRequest,
) (*url_shortener.NewShortUrlResponse, error) {
	if request.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url required")
	}
	fmt.Printf("url req %s", request.Url)
	shortUrl, err := s.usecase.CreateShortUrl(ctx, request.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &url_shortener.NewShortUrlResponse{ShortUrl: shortUrl}, nil
}

func (s *ServerApi) GetOriginalLink(
	ctx context.Context,
	request *url_shortener.GetOriginalUrlRequest,
) (*url_shortener.GetOriginalUrlResponse, error) {
	if len(request.ShortUrl) != 10 {
		return nil, status.Error(codes.InvalidArgument, "wrong short url len")
	}

	originalUrl, err := s.usecase.GetOriginalUrl(ctx, request.ShortUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &url_shortener.GetOriginalUrlResponse{OriginalUrl: originalUrl}, nil
}
