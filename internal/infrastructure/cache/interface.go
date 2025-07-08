package cache

import (
	"context"
	"time"
)

// Cache defines the interface for cache operations
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in the cache
	Set(ctx context.Context, key string, value string, expiration time.Duration) error

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Lock attempts to acquire a lock for the given key with a specified expiration duration.
	// Returns true if the lock is successfully acquired, otherwise false along with any error encountered.
	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)

	// Release releases a previously acquired lock for the given key in the cache.
	Release(ctx context.Context, key string) error

	// Close closes the cache connection
	Close() error
}
