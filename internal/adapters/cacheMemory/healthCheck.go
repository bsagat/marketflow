package cache

import (
	"context"
)

func (c *RedisCacheMemory) CheckHealth() error {
	_, err := c.Cache.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}
