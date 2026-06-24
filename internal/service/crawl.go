package service

import (
	"context"
	"tenome/internal/cerrors"
	"tenome/internal/worker"
)

type Crawler interface {
	Submit(ctx context.Context, url string) error
}

type CrawlService struct {
	jobs chan<- worker.CrawlJob
}

func NewCrawlService(jobs chan<- worker.CrawlJob) *CrawlService {
	return &CrawlService{jobs}
}

func (s *CrawlService) Submit(ctx context.Context, url string) error {
	select {
	case s.jobs <- worker.CrawlJob{URL: url}:
		return nil
	default:
		return cerrors.ErrQueueFull
	}
}
