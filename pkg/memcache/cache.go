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

	span, ctx := trace.New(context.Background(), "cache.CacheWorker")
	defer span.Close()

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.clear(ctx)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func RunPersistWorker(cache *MemCache) {

	span, ctx := trace.New(context.Background(), "cache.PersistWorker")
	defer span.Close()

	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.DB.Persist(ctx)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func (c *MemCache) clear(ctx context.Context) {

	span, ctx := trace.New(ctx, "cache.clear")
	defer span.Close()

	c.UserCache.data = make(map[string]*CachedUser)
	c.HeartCache.data = make(map[string]*CachedHeart)

}
