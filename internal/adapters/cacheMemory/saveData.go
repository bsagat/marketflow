package cache

import (
	"context"
	"encoding/json"
	"marketflow/internal/domain"
)

func (c *RedisCacheMemory) SaveAggregatedData(aggregatedData map[string]domain.ExchangeData) error {
	for key, value := range aggregatedData {
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if err := c.Cache.Set(context.Background(), key, jsonData, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}
