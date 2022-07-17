//go:generate mockgen -package=cache -destination=./cache_mock.go . Cache
package cache

import (
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"news-app/internal/domain"
)

// Cache is an interface for interacting with a caching layer
type Cache interface {
	GetArticlesFromCache(url string) ([]domain.Article, bool)
	AddArticlesToCache(url string, articles []domain.Article)
}

// cache is the internal representation of our cache
type cache struct {
	ttl   time.Duration
	clock clockwork.Clock

	mutex          sync.RWMutex
	feedToArticles map[string]cachedArticles
}

type cachedArticles struct {
	created  time.Time
	articles []domain.Article
}

// NewCache is a constructor for a Cache
// ttlDuration represents how long a cache entry should live
// tickerDuration represents how long between each cache evaluation
func NewCache(ttlDuration, tickerDuration time.Duration, clock clockwork.Clock) Cache {
	cache := &cache{
		ttl:            ttlDuration,
		feedToArticles: make(map[string]cachedArticles),
		clock:          clock,
	}

	ticker := clock.NewTicker(tickerDuration)
	go cache.cleanup(ticker)

	return cache
}

// GetArticlesFromCache will return true if cache was a hit, along with the articles it found. It will return false on a miss and a nil slice
func (c *cache) GetArticlesFromCache(url string) ([]domain.Article, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	v, ok := c.feedToArticles[url]
	if !ok {
		return nil, ok
	}

	return v.articles, ok
}

// AddArticlesToCache will add or overwrite article to cache
func (c *cache) AddArticlesToCache(url string, articles []domain.Article) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.feedToArticles[url] = cachedArticles{
		created:  c.clock.Now(),
		articles: articles,
	}
}

// cleanup will evaluate the cache every time the ttl has passed and delete any records that have existed longer than the ttl
func (c *cache) cleanup(ticker clockwork.Ticker) {
	defer ticker.Stop()

	for {
		select {
		case <-ticker.Chan():
			c.mutex.Lock()
			for k, v := range c.feedToArticles {
				// if the cache entry was created more than X ago
				if c.clock.Since(v.created) >= c.ttl {
					delete(c.feedToArticles, k)
				}
			}
			c.mutex.Unlock()
		}
	}
}
