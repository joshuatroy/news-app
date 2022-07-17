//go:generate mockgen -package=service -destination=./service_mock.go . Service

package service

import (
	"context"
	"fmt"
	"news-app/internal/cache"
	"sort"

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
	cache  cache.Cache
}

// NewService is a constructor for a Service
func NewService(parser parser.UniversalParser, cache cache.Cache) Service {
	return &service{
		cache:  cache,
		parser: parser,
	}
}

// GetArticles returns a list of articles given a feed URL
func (s service) GetArticles(ctx context.Context, feedURL string) ([]domain.Article, error) {
	articles, ok := s.cache.GetArticlesFromCache(feedURL)
	if !ok {
		feed, err := s.parser.Parse(ctx, feedURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse feed: %w", err)
		}

		// sort articles in descending order by published date
		sort.Slice(feed.Articles, func(i, j int) bool {
			return feed.Articles[i].Published.After(feed.Articles[j].Published)
		})

		articles = feed.Articles

		s.cache.AddArticlesToCache(feedURL, articles)
	}

	return articles, nil
}
