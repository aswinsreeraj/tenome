package service

import (
	"context"
	"log"
	"strings"
	"tenome/internal/cache"
	"tenome/internal/index"
	"tenome/internal/model"
	"tenome/internal/storage"
)

type SearchService struct {
	index   index.Index
	storage storage.Storage
	cache   cache.Cache
}

func NewSearchService(index index.Index, storage storage.Storage, cache cache.Cache) *SearchService {
	return &SearchService{index, storage, cache}
}

func (s *SearchService) Search(ctx context.Context, term string) ([]model.Page, error) {
	term = strings.ToLower(strings.TrimSpace(term))

	key := "search:" + term

	pages, found, err := s.cache.Get(ctx, key)
	if err != nil {
		found = false
	}
	if found {
		return pages, nil
	}

	ids, err := s.index.Search(ctx, term)

	if err != nil {
		return nil, err
	}

	pages, err = s.storage.GetPagesByIDs(ctx, ids)

	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, key, pages); err != nil {
		log.Println("Cache saving failed.", err)
	}

	return pages, nil

}
