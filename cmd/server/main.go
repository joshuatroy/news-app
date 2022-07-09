package main

import (
	"news-app/internal/parser"

	"github.com/mmcdole/gofeed"
)

var (
	//timeout on calls to rss feeds
	timeout = 10
)

func main() {
	_ = parser.NewParser(timeout, gofeed.NewParser())
}
