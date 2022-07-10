//go:generate mockgen -package=service -destination=./service_mock.go . Service

package service

import (
	"context"
	"fmt"

	"news-app/internal/domain"
	"news-app/internal/parser"
)

// Service interface represents the service layer function available
type Service interface {
	GetArticles(context.Context, string) ([]domain.Article, error)
}

// service is our internal representation of our service
type service struct {
	parser parser.UniversalParser
}

// NewService is a constructor for a Service
func NewService(parser parser.UniversalParser) Service {
	return &service{
		parser: parser,
	}
}

// GetArticles returns a list of articles given a feed URL
func (s service) GetArticles(ctx context.Context, feedURL string) ([]domain.Article, error) {
	feed, err := s.parser.Parse(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	return feed.Articles, nil
}
