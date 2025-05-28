package cache

import (
	"context"
	"encoding/json"
	"marketflow/internal/domain"
	"time"
)

var _ (domain.CacheMemory) = (*RedisCacheMemory)(nil)

func (c *RedisCacheMemory) SaveAggregatedData(aggregatedData map[string]domain.ExchangeData) error {
	for key, value := range aggregatedData {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err = c.Cache.Set(ctx, key, jsonData, 0).Err()
		cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

// Saves latest prices for every exchange every second
// key: Latest {Exchange} {Symbol}
func (c *RedisCacheMemory) SaveLatestData(latestData map[string]domain.Data) error {
	for key, value := range latestData {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err = c.Cache.Set(ctx, key, jsonData, 0).Err()
		cancel()
		if err != nil {
			return err
		}
	}

	return nil
}
