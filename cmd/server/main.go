package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"tenome/internal/api"
	"tenome/internal/cache"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/service"
	"tenome/internal/storage"
	"tenome/internal/worker"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()
	dbPath := os.Getenv("DB_PATH")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
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
