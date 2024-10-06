package redis_cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/Slava02/Involvio/bot/config"
	"github.com/Slava02/Involvio/bot/pkg/cache"
	tm "github.com/and3rson/telemux/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	redisConfigKey = "Redis"
)

func NewClient() *RedisClient {
	return new(RedisClient)
}

type RedisClient struct {
	Client *redis.Client
}

func (c *RedisClient) GetState(pk tm.PersistenceKey) string {
	//TODO implement me
	panic("implement me")
}

func (c *RedisClient) SetState(pk tm.PersistenceKey, state string) {
	//TODO implement me
	panic("implement me")
}

func (c *RedisClient) GetData(pk tm.PersistenceKey) tm.Data {
	//TODO implement me
	panic("implement me")
}

func (c *RedisClient) SetData(pk tm.PersistenceKey, data tm.Data) {
	//TODO implement me
	panic("implement me")
}

func (c *RedisClient) Configure(ctx context.Context, config *config.Config) {
	cfg := config.CacheConfig

	c.Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,
	})

	if err := c.Client.Ping(ctx).Err(); err != nil {
		logrus.Panicf("can't create redis Client, check Redis.json, err: %v", err)
	}

	logrus.Debugf("redis Client inited")
}

type Cache[T any] struct {
	*RedisClient
}

func New[T any](client *RedisClient) *Cache[T] {
	return &Cache[T]{
		RedisClient: client,
	}
}

func (c *Cache[T]) Set(ctx context.Context, key string, value T, ttl int) error {
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.Client.Set(
		ctx, key, data,
		time.Duration(ttl)*time.Hour,
	).Err(); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

func (c *Cache[T]) Get(ctx context.Context, key string) (T, error) {
	var defaultValue T

	data, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return defaultValue, cache.ErrKeyNotFound
		}

		return defaultValue, fmt.Errorf("failed to get value: %w", err)
	}

	var value T
	if err := jsoniter.Unmarshal([]byte(data), &value); err != nil {
		return defaultValue, fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return value, nil
}

func (c *Cache[T]) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.Client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get key ttl: %w", err)
	}

	return ttl, nil
}

func (c *Cache[T]) GetAll(ctx context.Context, pattern string) (map[string]T, error) {
	keys, err := c.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get all keys: %w", err)
	}

	result := make(map[string]T)
	for _, key := range keys {
		value, err := c.Get(ctx, key)
		if err == nil {
			result[key] = value
		}
	}
	return result, nil
}

func (c *Cache[T]) Update(ctx context.Context, key string, newValue T) error {
	exists, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check key exists: %w", err)
	}
	if exists == 0 {
		return cache.ErrKeyNotFound
	}

	data, err := jsoniter.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.Client.Set(
		ctx, key, data,
		redis.KeepTTL,
	).Err(); err != nil {
		return fmt.Errorf("failed to set newValue: %w", err)
	}

	return nil
}

func (c *Cache[T]) Delete(ctx context.Context, pattern string) error {
	keys, err := c.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get all keys: %w", err)
	}

	for _, key := range keys {
		if err := c.Client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("failed to delete key: %w", err)
		}
	}
	return nil
}

func (c *Cache[T]) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key exists: %w", err)
	}

	return count > 0, nil
}
