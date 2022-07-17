package cache

import (
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"news-app/internal/domain"
	"testing"
	"time"
)

func Test_Cache(t *testing.T) {
	const (
		someURL         = "some-url"
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
		someArticles       = []domain.Article{someArticle}
		someTTLDuration    = time.Duration(5) * time.Second
		someTickerDuration = time.Duration(1) * time.Second
	)
	t.Run("should return cache hit and articles", func(t *testing.T) {
		cache := NewCache(someTTLDuration, someTickerDuration, clockwork.NewRealClock())

		// add item to cache
		cache.AddArticlesToCache(someURL, someArticles)

		// get item from cache
		articles, ok := cache.GetArticlesFromCache(someURL)
		assert.True(t, ok)
		assert.Equal(t, articles, someArticles)
	})
	t.Run("should return cache miss and nil articles", func(t *testing.T) {
		cache := NewCache(someTTLDuration, someTickerDuration, clockwork.NewRealClock())

		// get item from cache
		articles, ok := cache.GetArticlesFromCache(someURL)
		assert.False(t, ok)
		assert.Nil(t, articles)
	})

	t.Run("should remove items from cache once ttl has passed", func(t *testing.T) {
		clock := clockwork.NewFakeClock()
		cache := NewCache(someTTLDuration, someTickerDuration, clock)

		// add item to cache
		cache.AddArticlesToCache(someURL, someArticles)
		articles, ok := cache.GetArticlesFromCache(someURL)
		assert.True(t, ok)
		assert.Equal(t, articles, someArticles)

		// advance clock and ticker so cache is cleaned up
		clock.Advance(someTTLDuration)
		clock.BlockUntil(1)

		// block to allow cache to be cleared before trying to access
		time.Sleep(10 * time.Millisecond)

		// fail to get item from cache
		articles, ok = cache.GetArticlesFromCache(someURL)
		assert.False(t, ok)
		assert.Nil(t, articles)
	})
}
