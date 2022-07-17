//go:generate mockgen -package=parser -destination=./parser_mock.go . InternalParser,UniversalParser

package parser

import (
	"context"
	"fmt"
	"news-app/internal/domain"
	"time"

	"github.com/mmcdole/gofeed"
)

// UniversalParser is an interface for parsing RSS feeds
type UniversalParser interface {
	Parse(ctx context.Context, url string) (domain.Feed, error)
}

// InternalParser an interface for mocking the gofeed package
// The gofeed package provides a parser that works with RSS, ATOM and JSON feeds.
type InternalParser interface {
	ParseURLWithContext(string, context.Context) (*gofeed.Feed, error)
}

// NewParser is a constructor for creating a parser.
func NewParser(timeout int, internalParser InternalParser) UniversalParser {
	return parser{
		timeout:        timeout,
		internalParser: internalParser,
	}
}

// parser is the internal representation of an RSS parser
type parser struct {
	timeout        int
	internalParser InternalParser
}

// Parse function will parse a feed from a FeedURL to a domain.Feed model
func (p parser) Parse(ctx context.Context, url string) (domain.Feed, error) {
	// Set up a timeout on network call
	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.timeout)*time.Second)
	defer cancel()

	// Call parser through Universal Parser interface so we can mock behaviour for testing
	feed, err := p.internalParser.ParseURLWithContext(url, ctx)
	if err != nil {
		return domain.Feed{}, fmt.Errorf("failed to parse url: %w", err)
	}

	return mapFeedToDomainModel(feed), nil
}

func mapFeedToDomainModel(f *gofeed.Feed) domain.Feed {
	if f != nil {
		var articles []domain.Article
		for _, item := range f.Items {
			articles = append(articles, mapItemToDomainModel(item))
		}

		return domain.Feed{
			Title:       f.Title,
			Description: f.Description,
			Articles:    articles,
		}
	}

	return domain.Feed{}
}

func mapItemToDomainModel(i *gofeed.Item) domain.Article {
	if i != nil {
		var published time.Time
		if i.PublishedParsed != nil {
			published = *i.PublishedParsed
		}

		return domain.Article{
			Title:       i.Title,
			Description: i.Description,
			Content:     i.Content,
			Image:       mapImageToDomainModel(i.Image),
			URL:         i.Link,
			Published:   published,
		}
	}

	return domain.Article{}
}

func mapImageToDomainModel(i *gofeed.Image) domain.Image {
	if i != nil {
		return domain.Image{
			Title: i.Title,
			URL:   i.URL,
		}
	}

	return domain.Image{}
}
