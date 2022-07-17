package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-app/internal/cache"
	"news-app/internal/domain"
	"news-app/internal/parser"
	"testing"
)

func Test_service_GetArticles(t *testing.T) {
	const (
		someURL         = "some-url"
		someFeedURL     = "some-feed-url"
		someTitle       = "some-title"
		someDescription = "some-description"
		someContent     = "some-content"
	)
	var (
		someImage = domain.Image{
			URL:   someURL,
			Title: someTitle,
		}
		someArticle = domain.Article{
			Title:       someTitle,
			Description: someDescription,
			Content:     someContent,
			Image:       someImage,
		}
		someArticles = []domain.Article{someArticle}
		someFeed     = domain.Feed{
			Title:       someTitle,
			Description: someDescription,
			Articles:    someArticles,
		}
	)
	t.Run("should return a list of articles from cache", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockParser := parser.NewMockUniversalParser(ctrl)
		mockCache := cache.NewMockCache(ctrl)
		service := NewService(mockParser, mockCache)

		mockCache.EXPECT().GetArticlesFromCache(someFeedURL).Return(someArticles, true)

		articles, err := service.GetArticles(context.Background(), someFeedURL)
		assert.NoError(t, err)
		assert.Equal(t, someFeed.Articles, articles)
	})
	t.Run("should parse a list of articles and add to cache", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockParser := parser.NewMockUniversalParser(ctrl)
		mockCache := cache.NewMockCache(ctrl)
		service := NewService(mockParser, mockCache)

		mockCache.EXPECT().GetArticlesFromCache(someFeedURL).Return(nil, false)
		mockParser.EXPECT().Parse(gomock.Any(), someFeedURL).Return(someFeed, nil)
		mockCache.EXPECT().AddArticlesToCache(someFeedURL, someArticles)

		articles, err := service.GetArticles(context.Background(), someFeedURL)
		assert.NoError(t, err)
		assert.Equal(t, someFeed.Articles, articles)
	})
	t.Run("should return an error if we fail to get a list of articles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockParser := parser.NewMockUniversalParser(ctrl)
		mockCache := cache.NewMockCache(ctrl)
		service := NewService(mockParser, mockCache)

		mockCache.EXPECT().GetArticlesFromCache(someFeedURL).Return(nil, false)
		mockParser.EXPECT().Parse(gomock.Any(), someFeedURL).Return(domain.Feed{}, assert.AnError)

		articles, err := service.GetArticles(context.Background(), someFeedURL)
		assert.Error(t, err)
		assert.Empty(t, articles)
	})
}
