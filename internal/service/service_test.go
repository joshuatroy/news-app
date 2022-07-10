package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		someFeed = domain.Feed{
			Title:       someTitle,
			Description: someDescription,
			Articles:    []domain.Article{someArticle},
		}
	)
	t.Run("should return a list of articles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockParser := parser.NewMockUniversalParser(ctrl)
		service := NewService(mockParser)

		mockParser.EXPECT().Parse(gomock.Any(), someFeedURL).Return(someFeed, nil)

		articles, err := service.GetArticles(context.Background(), someFeedURL)
		assert.NoError(t, err)
		assert.Equal(t, someFeed.Articles, articles)
	})

	t.Run("should return an error if we fail to get a list of articles", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockParser := parser.NewMockUniversalParser(ctrl)
		service := NewService(mockParser)

		mockParser.EXPECT().Parse(gomock.Any(), someFeedURL).Return(domain.Feed{}, assert.AnError)

		articles, err := service.GetArticles(context.Background(), someFeedURL)
		assert.Error(t, err)
		assert.Empty(t, articles)
	})
}
