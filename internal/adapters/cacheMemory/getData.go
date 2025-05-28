package cache

import (
	"context"
	"encoding/json"
	"marketflow/internal/domain"
	"time"
)

// Latest Data fetching from Redis cache memory
//
// Argument parameters:
//   - Send only valid data
//   - Key structure : "[exchangeNum] [symbol]"
func (c *RedisCacheMemory) GetLatestData(exchange, symbol string) (domain.Data, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := "latest " + exchange + " " + symbol
	res, err := c.Cache.Get(ctx, key).Result()
	if err != nil {
		return domain.Data{}, err
	}

	raw := domain.Data{}
	if err := json.Unmarshal([]byte(res), &raw); err != nil {
		return domain.Data{}, err
	}
	return raw, nil
}
