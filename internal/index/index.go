package index

import (
	"context"
	"strings"
	"sync"
	"tenome/internal/model"
)

type Index interface {
	Add(ctx context.Context, page model.Page) error
	Search(ctx context.Context, term string) ([]int64, error)
}

type InvertedIndex struct {
	index map[string]map[int64]struct{}
	mu    sync.Mutex
}

func New() *InvertedIndex {
	index := make(map[string]map[int64]struct{})

	return &InvertedIndex{index: index}
}

func (i *InvertedIndex) Add(ctx context.Context, page model.Page) error {
	tokens := tokenize(page)

	// Mutex for locking the index
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, token := range tokens {
		if _, exists := i.index[token]; !exists {
			i.index[token] = make(map[int64]struct{})
		}
		i.index[token][page.ID] = struct{}{}
	}

	return nil
}

func tokenize(page model.Page) []string {
	combined := page.Title + " " + page.Content
	combined = strings.ToLower(combined)
	tokens := strings.Fields(combined)
	dedup := make(map[string]struct{})
	dedupToken := []string{}
	for _, str := range tokens {
		token := strings.Trim(str, ".,!?;:\"'(){}[]")
		if _, exists := dedup[token]; exists {
			continue
		}
		dedup[token] = struct{}{}
		dedupToken = append(dedupToken, token)

	}
	return dedupToken
}

func (i *InvertedIndex) Search(ctx context.Context, term string) ([]int64, error) {
	term = strings.ToLower(term)
	list, exists := i.index[term]
	if !exists {
		return []int64{}, nil
	}
	index := []int64{}

	for key := range list {
		index = append(index, key)
	}

	return index, nil
}
