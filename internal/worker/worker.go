package worker

import (
	"context"
	"fmt"
	"tenome/internal/crawler"
	"tenome/internal/index"
	"tenome/internal/storage"
)

type Worker struct {
	id      int
	crawler crawler.Crawler
	storage storage.Storage
	index   index.Index
}

func New(id int, crawler crawler.Crawler, storage storage.Storage, index index.Index) Worker {
	return Worker{id, crawler, storage, index}
}

func (w *Worker) Start(ctx context.Context, jobs <-chan CrawlJob) {

	for job := range jobs {
		fmt.Printf("Worker #%d pocessing %s\n", w.id, job.URL)
		if err := w.Process(ctx, job); err != nil {
			continue
		}
	}

}

func (w *Worker) Process(ctx context.Context, job CrawlJob) error {
	page, err := w.crawler.Crawl(ctx, job.URL)
	if err != nil {
		return err
	}
	id, err := w.storage.SavePage(ctx, page)
	if err != nil {
		return err
	}
	page.ID = id

	return w.index.Add(ctx, page)
}
