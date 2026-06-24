package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"tenome/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, pages []model.Page) error {
	serial, err := json.Marshal(pages)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}
	err = r.client.Set(ctx, key, serial, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("cache db failure: %w", err)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]model.Page, bool, error) {
	bypages, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("cache failure: %w", err)
	}

	var pages []model.Page

	err = json.Unmarshal([]byte(bypages), &pages)

	if err != nil {
		return nil, true, fmt.Errorf("Unmarshal error: %w", err)
	}

	return pages, true, nil

}
