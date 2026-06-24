package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/model"
	"tenome/internal/storage"

	_ "modernc.org/sqlite"
)

func main() {

	ctx := context.Background()
	db, err := sql.Open("sqlite", "crawler.db")
	if err != nil {
		panic(err)
	}

	storage := storage.New(db)
	err = storage.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	pagex := model.Page{
		URL:     "https://go.dev",
		Title:   "Go",
		Content: "Go language",
	}

	page1 := model.Page{URL: "https://aswingo.dev", Title: "Golang Fundamentals", Content: "Go is good for concurrency"}
	page2 := model.Page{URL: "https://go.aswin", Title: "Go!", Content: "Channels"}
	storage.SavePage(ctx, pagex)

	storage.SavePage(ctx, page1)
	storage.SavePage(ctx, page2)
	pages, err := storage.GetPagesByIDs(ctx, []int64{1, 2, 3})

	fmt.Println(pages)

	idx := index.New()
	idx.Add(ctx, page1)
	idx.Add(ctx, page2)

	inds, _ := idx.Search(ctx, "GO")
	fmt.Println(inds)

	client := http.Client{}

	crawly := crawler.New(&client)
	page, _ := crawly.Crawl(ctx, "https://go.dev")

	fmt.Println(page.Title, page.URL, page.ID)
	fmt.Println(len(page.Content))
}
