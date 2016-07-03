package memcache

import (
	"github.com/romainmenke/hearts/pkg/fakedb"
	"golang.org/x/net/context"
	"limbo.services/trace"
)

type HeartCache struct {
	data map[string]*CachedHeart
}

type CachedHeart struct {
	Heart *fakedb.Heart
	Etag  int
}

func (c *MemCache) LoadHeart(ctx context.Context, domain string, owner string, repo string) (*CachedHeart, error) {

	span, ctx := trace.New(ctx, "cache.LoadHeart")
	defer span.Close()

	heart := &fakedb.Heart{
		Domain: domain,
		Owner:  owner,
		Repo:   repo,
	}

	cached, exists := c.HeartCache.data[heart.FullPath()]
	if exists {
		return cached, nil
	}

	heart, err := c.DB.LoadHeart(ctx, domain, owner, repo)
	if err != nil {
		return nil, err
	}

	cached = &CachedHeart{
		Heart: heart,
		Etag:  heart.Hash(),
	}

	c.HeartCache.data[heart.FullPath()] = cached

	return cached, nil
}

func (c *MemCache) SaveHeart(ctx context.Context, heart *fakedb.Heart) error {

	span, ctx := trace.New(ctx, "cache.SaveUser")
	defer span.Close()

	cached, exists := c.HeartCache.data[heart.FullPath()]
	if exists && cached.Heart == heart {
		return nil
	}

	err := c.DB.SaveObject(ctx, heart)
	if err != nil {
		return err
	}

	c.HeartCache.data[heart.FullPath()] = &CachedHeart{
		Heart: heart,
		Etag:  heart.Hash(),
	}

	return nil
}
