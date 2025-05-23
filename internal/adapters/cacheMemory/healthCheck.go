package cache

import (
	"context"
	"fmt"
)

func (c *RedisCacheMemory) CheckHealth() error {
	msg, err := c.Cache.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	fmt.Println(msg)
	return nil
}
