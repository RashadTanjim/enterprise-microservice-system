package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis cache configuration.
type Config struct {
	Enabled    bool
	Host       string
	Port       string
	Password   string
	DB         int
	DefaultTTL time.Duration
}

// Cache provides JSON-based Redis caching with optional disablement.
type Cache struct {
	client     *redis.Client
	enabled    bool
	defaultTTL time.Duration
	prefix     string
}

// NewRedisCache creates a cache client and verifies connectivity when enabled.
func NewRedisCache(cfg Config, prefix string) (*Cache, error) {
	cache := &Cache{
		enabled:    cfg.Enabled,
		defaultTTL: cfg.DefaultTTL,
		prefix:     prefix,
	}

	if !cfg.Enabled {
		return cache, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		cache.enabled = false
		return cache, err
	}

	cache.client = client
	return cache, nil
}

// Enabled returns whether caching is active.
func (c *Cache) Enabled() bool {
	return c != nil && c.enabled && c.client != nil
}

// GetJSON retrieves a cached JSON value into dest. Returns false when missing/disabled.
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	if !c.Enabled() {
		return false, nil
	}

	value, err := c.client.Get(ctx, c.key(key)).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return false, err
	}

	return true, nil
}

// SetJSON stores a value as JSON. When ttl <= 0, default TTL is used.
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !c.Enabled() {
		return nil
	}

	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.key(key), payload, ttl).Err()
}

// Delete removes keys from cache. No-op when disabled.
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if !c.Enabled() || len(keys) == 0 {
		return nil
	}

	prefixed := make([]string, 0, len(keys))
	for _, key := range keys {
		prefixed = append(prefixed, c.key(key))
	}

	return c.client.Del(ctx, prefixed...).Err()
}

func (c *Cache) key(key string) string {
	if c.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.prefix, key)
}
