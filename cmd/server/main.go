package main

import (
	"context"
	"fmt"
	"net/http"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/model"
)

func main() {
	ctx := context.Background()
	page1 := model.Page{ID: 1, Title: "Golang Fundamentals", Content: "Go is good for concurrency"}
	page2 := model.Page{ID: 2, Title: "Go!", Content: "Channels"}
	idx := index.New()
	idx.Add(ctx, page1)
	idx.Add(ctx, page2)

	inds, _ := idx.Search(ctx, "GO")
	fmt.Println(inds)

	client := http.Client{}

	crawly := crawler.New(&client)
	_, err := crawly.Crawl(ctx, "https://crawler-test.com/status_codes/status_404")

	fmt.Println(err)
}
