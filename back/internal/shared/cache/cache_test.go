package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newStore(t *testing.T) (*Store, *miniredis.Miniredis) {
	t.Helper()
	server := miniredis.RunT(t)
	return New(redis.NewClient(&redis.Options{Addr: server.Addr()})), server
}

var policy = Policy{FreshFor: time.Hour, KeepFor: 24 * time.Hour}

func TestGetOrRefreshServesFreshValueWithoutRefreshing(t *testing.T) {
	store, _ := newStore(t)
	calls := 0
	refresh := func(context.Context) (string, error) { calls++; return "value", nil }

	for range 2 {
		value, err := GetOrRefresh(context.Background(), store, "key", policy, refresh)
		if err != nil || value != "value" {
			t.Fatalf("GetOrRefresh = %q, %v", value, err)
		}
	}
	if calls != 1 {
		t.Fatalf("refresh called %d times, want 1", calls)
	}
}

func TestGetOrRefreshServesStaleValueWhenRefreshFails(t *testing.T) {
	store, _ := newStore(t)
	expired := Policy{FreshFor: -time.Second, KeepFor: 24 * time.Hour}
	if _, err := GetOrRefresh(context.Background(), store, "key", expired, func(context.Context) (string, error) { return "old", nil }); err != nil {
		t.Fatal(err)
	}

	value, err := GetOrRefresh(context.Background(), store, "key", policy, func(context.Context) (string, error) {
		return "", errors.New("source down")
	})

	if err != nil || value != "old" {
		t.Fatalf("GetOrRefresh = %q, %v; want the stale value", value, err)
	}
}

func TestGetOrRefreshReportsErrorWithoutStaleValue(t *testing.T) {
	store, _ := newStore(t)
	failure := errors.New("source down")

	_, err := GetOrRefresh(context.Background(), store, "key", policy, func(context.Context) (string, error) { return "", failure })

	if !errors.Is(err, failure) {
		t.Fatalf("error = %v, want %v", err, failure)
	}
}

func TestGetOrRefreshWorksWithoutRedis(t *testing.T) {
	store, server := newStore(t)
	server.Close()

	value, err := GetOrRefresh(context.Background(), store, "key", policy, func(context.Context) (string, error) { return "value", nil })

	if err != nil || value != "value" {
		t.Fatalf("GetOrRefresh = %q, %v", value, err)
	}
}
