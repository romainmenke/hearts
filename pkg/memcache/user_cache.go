package memcache

import (
	"github.com/romainmenke/hearts/pkg/fakedb"
	"golang.org/x/net/context"
	"limbo.services/trace"
)

type UserCache struct {
	data map[string]*CachedUser
}

type CachedUser struct {
	User *fakedb.User
	Etag int
}

func (c *MemCache) LoadUser(ctx context.Context, domain string, name string) (*CachedUser, error) {

	span, ctx := trace.New(ctx, "cache.LoadUser")
	defer span.Close()

	user := &fakedb.User{
		Domain: domain,
		Name:   name,
	}
	cached, exists := c.UserCache.data[user.FullPath()]
	if exists {
		return cached, nil
	}

	user, err := c.DB.LoadUser(ctx, domain, name)
	if err != nil {
		return nil, err
	}

	cached = &CachedUser{
		User: user,
		Etag: user.Hash(),
	}

	c.UserCache.data[user.FullPath()] = cached

	return cached, nil
}

func (c *MemCache) SaveUser(ctx context.Context, user *fakedb.User) error {

	span, ctx := trace.New(ctx, "cache.SaveUser")
	defer span.Close()

	cached, exists := c.UserCache.data[user.FullPath()]
	if exists && cached.Etag == user.Hash() {
		return nil
	}

	err := c.DB.SaveObject(ctx, user)
	if err != nil {
		return err
	}

	c.UserCache.data[user.FullPath()] = &CachedUser{
		User: user,
		Etag: user.Hash(),
	}

	return nil
}
