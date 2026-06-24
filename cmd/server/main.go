package main

import (
	"context"
	"database/sql"
	"net/http"
	"sync"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/storage"
	"tenome/internal/worker"

	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

func main() {

	ctx := context.Background()
	var wg sync.WaitGroup
	db, err := sql.Open("sqlite", "crawler.db")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

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

		wg.Go(func() {
			w.Start(ctx, jobs)
		})
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

	wg.Wait()

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
