package parser

import (
	"context"
	"news-app/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func Test_parser_Parse(t *testing.T) {
	var (
		someURL         = "https://some-url.com"
		someTitle       = "someTitle"
		someDescription = "someDescription"
		someContent     = "someContent"
		someImage       = gofeed.Image{
			URL:   someURL,
			Title: someTitle,
		}
		someItem = gofeed.Item{
			Title:       someTitle,
			Description: someDescription,
			Content:     someContent,
			Image:       &someImage,
		}
		someFeed = gofeed.Feed{
			Title:       someTitle,
			Description: someDescription,
			Items:       []*gofeed.Item{&someItem},
		}
	)

	t.Run("parser should return a feed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockInternalParser := NewMockInternalParser(ctrl)

		parser := NewParser(10, mockInternalParser)

		mockInternalParser.EXPECT().ParseURLWithContext(someURL, gomock.Any()).Return(&someFeed, nil)

		feed, err := parser.Parse(context.Background(), someURL)
		assert.NoError(t, err)

		expected := domain.Feed{
			Title:       someTitle,
			Description: someDescription,
			Articles: []domain.Article{
				{
					Title:       someTitle,
					Description: someDescription,
					Content:     someContent,
					Image: domain.Image{
						Title: someTitle,
						URL:   someURL,
					},
				},
			},
		}
		assert.Equal(t, expected, feed)
	})
	t.Run("parser should return an error if we fail to parse FeedURL", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockInternalParser := NewMockInternalParser(ctrl)

		parser := NewParser(10, mockInternalParser)

		mockInternalParser.EXPECT().ParseURLWithContext(someURL, gomock.Any()).Return(nil, assert.AnError)

		feed, err := parser.Parse(context.Background(), someURL)
		assert.Error(t, err)
		assert.Empty(t, feed)
	})
}
