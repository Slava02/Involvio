package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Cache[T any] interface {
	Set(ctx context.Context, key string, value T, ttl int) error
	Get(ctx context.Context, key string) (T, error)
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	GetAll(ctx context.Context, pattern string) (map[string]T, error)
	Update(ctx context.Context, key string, newValue T) error
	Delete(ctx context.Context, pattern string) error
	Exists(ctx context.Context, pattern string) (bool, error)
}
