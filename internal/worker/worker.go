package worker

import (
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/storage"
)

type Worker struct {
	crawler crawler.Crawler
	storage storage.Storage
	index   index.Index
}

type CrawlJob struct {
	URL string
}
