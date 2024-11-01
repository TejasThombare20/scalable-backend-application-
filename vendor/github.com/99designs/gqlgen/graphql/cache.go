package graphql

import "context"

// Cache is a shared store for APQ and query AST caching
type Cache[T any] interface {
	// Get looks up a key's value from the cache.
	Get(ctx context.Context, key string) (value T, ok bool)

	// Add adds a value to the cache.
	Add(ctx context.Context, key string, value T)
}

// MapCache is the simplest implementation of a cache, because it can not evict it should only be used in tests
type MapCache[T any] map[string]T

// Get looks up a key's value from the cache.
func (m MapCache[T]) Get(_ context.Context, key string) (value T, ok bool) {
	v, ok := m[key]
	return v, ok
}

// Add adds a value to the cache.
func (m MapCache[T]) Add(_ context.Context, key string, value T) { m[key] = value }

type NoCache[T any, T2 *T] struct{}

var _ Cache[*string] = (*NoCache[string, *string])(nil)

func (n NoCache[T, T2]) Get(_ context.Context, _ string) (value T2, ok bool) { return nil, false }
func (n NoCache[T, T2]) Add(_ context.Context, _ string, _ T2)               {}
