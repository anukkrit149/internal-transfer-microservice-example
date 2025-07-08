package cache

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"internal-transfer-microservice/internal/config"
)

const LockPrefix = "lock:"

// RedisCache implements Cache interface
type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// Use SetNX (SET if Not exists) for distributed locking
	// This will set the key only if it doesn't already exist, if exist it will fail
	success, err := r.client.SetNX(ctx, LockPrefix+key, "locked", expiration).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *RedisCache) Release(ctx context.Context, key string) error {
	// Release the lock by deleting the key
	return r.Delete(ctx, key)
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(cfg *config.Config) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddress(),
		Password: cfg.GetRedisPassword(),
		DB:       cfg.GetRedisDB(),
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis cache")

	return &RedisCache{
		client: client,
	}, nil
}

// Get retrieves a value from the cache
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a value in the cache
func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Delete removes a value from the cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis client
func (r *RedisCache) Close() error {
	return r.client.Close()
}
