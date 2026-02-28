package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	prefix string
}

func New(client *redis.Client, prefix string) *Cache {
	if client == nil {
		return nil
	}
	return &Cache{client: client, prefix: prefix}
}

func (c *Cache) Key(key string) string {
	if c == nil {
		return ""
	}
	return c.prefix + key
}

func (c *Cache) GetString(ctx context.Context, key string) (string, bool, error) {
	if c == nil {
		return "", false, nil
	}
	val, err := c.client.Get(ctx, c.Key(key)).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) SetString(ctx context.Context, key string, value string, ttl time.Duration) error {
	if c == nil {
		return nil
	}
	return c.client.Set(ctx, c.Key(key), value, ttl).Err()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	if c == nil {
		return nil
	}
	if len(keys) == 0 {
		return nil
	}
	expanded := make([]string, 0, len(keys))
	for _, k := range keys {
		expanded = append(expanded, c.Key(k))
	}
	return c.client.Del(ctx, expanded...).Err()
}

func (c *Cache) DelPrefix(ctx context.Context, prefix string) error {
	if c == nil {
		return nil
	}

	pattern := c.Key(prefix) + "*"
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(ctx, cursor, pattern, 200).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}
