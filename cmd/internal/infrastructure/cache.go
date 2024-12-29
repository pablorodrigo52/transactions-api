package infrastructure

import (
	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	Cache *ristretto.Cache
}

func NewCache() *Cache {

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic("Failed to create cache " + err.Error())
	}

	return &Cache{
		Cache: cache,
	}
}
