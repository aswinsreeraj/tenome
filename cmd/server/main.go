package main

import (
	"context"
	"database/sql"
	"net/http"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/storage"
	"tenome/internal/worker"
	"time"

	_ "modernc.org/sqlite"
)

func main() {

	ctx := context.Background()
	db, err := sql.Open("sqlite", "crawler.db")
	if err != nil {
		panic(err)
	}

	storage := storage.New(db)
	idx := index.New()
	client := http.Client{}
	crawly := crawler.New(&client)

	err = storage.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	jobs := make(chan worker.CrawlJob, 100)
	for i := range 4 {
		w := worker.New(i, crawly, storage, idx)
		go w.Start(ctx, jobs)
	}

	urls := []string{
		"https://go.dev",
		"https://golang.org",
		"https://pkg.go.dev",
		"https://go.dev/doc",
	}

	for _, url := range urls {
		jobs <- worker.CrawlJob{URL: url}
	}

	close(jobs)

	time.Sleep(2 * time.Second)

	// pagex := model.Page{
	// 	URL:     "https://go.dev",
	// 	Title:   "Go",
	// 	Content: "Go language",
	// }

	// page1 := model.Page{URL: "https://aswingo.dev", Title: "Golang Fundamentals", Content: "Go is good for concurrency"}
	// page2 := model.Page{URL: "https://go.aswin", Title: "Go!", Content: "Channels"}
	// storage.SavePage(ctx, pagex)

	// storage.SavePage(ctx, page1)
	// storage.SavePage(ctx, page2)
	// pages, err := storage.GetPagesByIDs(ctx, []int64{1, 2, 3})

	// idx.Add(ctx, page1)
	// idx.Add(ctx, page2)

	// inds, _ := idx.Search(ctx, "GO")
	// fmt.Println(inds)

	// page, _ := crawly.Crawl(ctx, "https://go.dev")

	// fmt.Println(page.Title, page.URL, page.ID)
	// fmt.Println(len(page.Content))
}
