//go:generate mockgen -package=service -destination=./service_mock.go . Service

package service

import (
	"context"
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

func (s service) GetArticles(context.Context, string) ([]domain.Article, error) {
	//TODO: implement me
	return nil, nil
}
