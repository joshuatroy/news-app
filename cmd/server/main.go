package main

import (
	"log"
	netHTTP "net/http"
	"time"

	"news-app/internal/parser"
	"news-app/internal/service"
	"news-app/internal/transport/http"

	"github.com/mmcdole/gofeed"
)

var (
	//timeout on calls to rss feeds
	timeout = 10
)

func main() {
	universalParser := parser.NewParser(timeout, gofeed.NewParser())
	svc := service.NewService(universalParser)
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
