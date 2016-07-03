package memcache

import (
	"time"

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
	}
}

func RunCacheWorker(cache *MemCache) {

	ticker := time.NewTicker(2 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.clear()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func RunPersistWorker(cache *MemCache) {

	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cache.DB.Persist(context.Background())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func (c *MemCache) clear() {

	c.UserCache.data = make(map[string]*CachedUser)
	c.HeartCache.data = make(map[string]*CachedHeart)

}
