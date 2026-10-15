// Package cache stores JSON values in Redis with a freshness date, so that a stale
// value can still be served when its refresh fails.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store is a Redis-backed cache of JSON values.
type Store struct {
	client *redis.Client
}

// New returns a store using the given Redis client.
func New(client *redis.Client) *Store {
	return &Store{client: client}
}

type entry struct {
	FreshUntil time.Time       `json:"freshUntil"`
	Value      json.RawMessage `json:"value"`
}

// Policy tells how long a value is fresh and how long a stale copy is kept.
type Policy struct {
	FreshFor time.Duration
	KeepFor  time.Duration
}

// GetOrRefresh returns the cached value while it is fresh. Otherwise it calls refresh
// and stores the result; when refresh fails, the stale value is returned if one exists.
// An unavailable Redis only disables caching: refresh is called directly.
func GetOrRefresh[T any](ctx context.Context, store *Store, key string, policy Policy, refresh func(context.Context) (T, error)) (T, error) {
	cached, found := store.read(ctx, key)
	var stale T
	hasStale := found && json.Unmarshal(cached.Value, &stale) == nil
	if hasStale && time.Now().Before(cached.FreshUntil) {
		return stale, nil
	}

	value, err := refresh(ctx)
	if err != nil && hasStale {
		slog.WarnContext(ctx, "serving stale cache entry", "key", key, "error", err)
		return stale, nil
	}
	if err != nil {
		return value, err
	}
	store.write(ctx, key, value, policy)
	return value, nil
}

func (s *Store) read(ctx context.Context, key string) (entry, bool) {
	raw, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return entry{}, false
	}
	if err != nil {
		slog.WarnContext(ctx, "cache read failed", "key", key, "error", err)
		return entry{}, false
	}
	var cached entry
	if err := json.Unmarshal(raw, &cached); err != nil {
		return entry{}, false
	}
	return cached, true
}

func (s *Store) write(ctx context.Context, key string, value any, policy Policy) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return
	}
	raw, err := json.Marshal(entry{FreshUntil: time.Now().Add(policy.FreshFor), Value: encoded})
	if err != nil {
		return
	}
	if err := s.client.Set(ctx, key, raw, policy.KeepFor).Err(); err != nil {
		slog.WarnContext(ctx, "cache write failed", "key", key, "error", err)
	}
}
