package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"tenome/internal/api"
	"tenome/internal/cache"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/service"
	"tenome/internal/storage"
	"tenome/internal/worker"

	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

func main() {

	ctx := context.Background()

	db, err := sql.Open("sqlite", "crawler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

	storage := storage.New(db)
	err = storage.Migrate(ctx)
	if err != nil {
		panic(err)
	}
	idx := index.New()

	pages, err := storage.GetAllPages(ctx)
	if err != nil {
		panic(err)
	}
	for _, page := range pages {
		idx.Add(ctx, page)
	}

	client := http.Client{}
	crawly := crawler.New(&client)
	cache := cache.NewCache(rdb)
	jobs := make(chan worker.CrawlJob, 100)
	searchService := service.NewSearchService(idx, storage, cache)
	crawlService := service.NewCrawlService(jobs)
	handler := api.NewHandler(searchService, crawlService)

	mux := api.NewRouter(handler)

	for i := range 4 {
		w := worker.New(i, crawly, storage, idx)
		go w.Start(ctx, jobs)
	}

	log.Println("server starting on :8050")
	log.Fatal(http.ListenAndServe(":8050", mux))

}
