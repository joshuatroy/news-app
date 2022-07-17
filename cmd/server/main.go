package main

import (
	"log"
	netHTTP "net/http"
	"time"

	"news-app/internal/cache"
	"news-app/internal/parser"
	"news-app/internal/service"
	"news-app/internal/transport/http"

	"github.com/jonboulle/clockwork"
	"github.com/mmcdole/gofeed"
)

var (
	//timeout on calls to rss feeds
	timeout = 10
	//ttlDuration is the time to live for cache entries
	ttlDuration = 5 * time.Minute
	//tickerDuration is the time between each cache evaluation
	tickerDuration = 1 * time.Minute
)

func main() {
	internalCache := cache.NewCache(
		ttlDuration,
		tickerDuration,
		clockwork.NewRealClock(),
	)

	universalParser := parser.NewParser(
		timeout,
		gofeed.NewParser(),
	)

	svc := service.NewService(
		universalParser,
		internalCache,
	)

	handler := http.NewHandler(svc)
	handler.ApplyRoutes()

	server := netHTTP.Server{
		Handler:      handler,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Duration(timeout) * time.Second,
		ReadTimeout:  time.Duration(timeout) * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
