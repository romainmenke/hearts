package memcache

import (
	"time"

	"limbo.services/trace"

	"golang.org/x/net/context"

	"github.com/romainmenke/hearts/pkg/fakedb"
)

type MemCache struct {
	DB         *fakedb.FakeDB
	UserCache  *UserCache
	HeartCache *HeartCache
}

func New(db *fakedb.FakeDB) *MemCache {
	return &MemCache{
		DB: db,
		UserCache: &UserCache{
			data: make(map[string]*CachedUser),
		},
		HeartCache: &HeartCache{
			data: make(map[string]*CachedHeart),
		},
	}
}

func RunCacheWorker(cache *MemCache) {

	span, _ := trace.New(context.Background(), "server.cache.CacheWorker")
	defer span.Close()

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.clean(context.Background())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func (c *MemCache) clean(ctx context.Context) {

	span, ctx := trace.New(ctx, "server.cache.clean")
	defer span.Close()

	now := time.Now()

	for key, value := range c.HeartCache.data {
		diff := now.Sub(value.Time)
		if diff > 5*time.Minute {
			delete(c.HeartCache.data, key)
		}
	}

	for key, value := range c.UserCache.data {
		diff := now.Sub(value.Time)
		if diff > 5*time.Minute {
			delete(c.UserCache.data, key)
		}
	}

}
