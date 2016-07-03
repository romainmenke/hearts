package main

import (
	"github.com/pborman/uuid"
	"github.com/romainmenke/hearts/pkg/fakedb"
	"github.com/romainmenke/hearts/pkg/memcache"
	"golang.org/x/net/context"
	"limbo.services/trace"
)

func update(ctx context.Context, cache *memcache.MemCache, message *incomingMessage) error {

	span, ctx := trace.New(ctx, "server.hearts.update")
	defer span.Close()

	heart := loadHeart(ctx, cache, message)
	user := loadUser(ctx, cache, message)

	applyChanges(ctx, message, heart, user)

	err := saveHeart(ctx, cache, heart)
	if err != nil {
		return span.Error(err)
	}

	err = saveUser(ctx, cache, user)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func loadHeart(ctx context.Context, cache *memcache.MemCache, message *incomingMessage) *fakedb.Heart {

	span, ctx := trace.New(ctx, "server.hearts.loadHeart")
	defer span.Close()

	heartCache, err := cache.LoadHeart(ctx, message.repo.domain, message.repo.owner, message.repo.name)
	if err != nil {
		span.Error(err)

		pass := message.result.pass
		newHeart := &fakedb.Heart{
			ID:        uuid.New(),
			Count:     3,
			LastBuild: pass,
			Domain:    message.repo.domain,
			Owner:     message.repo.owner,
			Repo:      message.repo.name,
		}
		return newHeart
	}

	return heartCache.Heart
}

func loadUser(ctx context.Context, cache *memcache.MemCache, message *incomingMessage) *fakedb.User {

	span, ctx := trace.New(ctx, "server.hearts.loadUser")
	defer span.Close()

	userCache, err := cache.LoadUser(ctx, message.repo.domain, message.user.name)
	if err != nil {
		span.Error(err)

		newUser := &fakedb.User{
			ID:     uuid.New(),
			Domain: message.repo.domain,
			Name:   message.user.name,
			Level:  0,
			Exp:    0,
			Streak: 0,
			Deaths: 0,
		}
		return newUser
	}

	return userCache.User
}

func saveHeart(ctx context.Context, cache *memcache.MemCache, heart *fakedb.Heart) error {

	span, ctx := trace.New(ctx, "server.hearts.saveHeart")
	defer span.Close()

	err := cache.SaveHeart(ctx, heart)
	if err != nil {
		return span.Error(err)
	}
	return nil

}

func saveUser(ctx context.Context, cache *memcache.MemCache, user *fakedb.User) error {

	span, ctx := trace.New(ctx, "server.hearts.saveUser")
	defer span.Close()

	err := cache.SaveUser(ctx, user)
	if err != nil {
		return span.Error(err)
	}
	return nil

}

func applyChanges(ctx context.Context, message *incomingMessage, heart *fakedb.Heart, user *fakedb.User) {

	span, ctx := trace.New(ctx, "server.hearts.applyChanges")
	defer span.Close()

	if heart.ID == "" {
		heart.ID = uuid.New()
	}

	if user.ID == "" {
		user.ID = uuid.New()
	}

	kill := false
	save := false

	if message.result.pass == true && heart.LastBuild == false && heart.Count == 1 && heart.LastBuilderID != user.ID {
		save = true
	}

	if message.result.pass == false && heart.Count == 1 {
		kill = true
	}

	// HEART
	if message.result.pass == true && heart.LastBuild == false && heart.Count != 0 {
		heart.Count++
		if heart.Count > 3 {
			heart.Count = 3
		}
	} else if message.result.pass == true && heart.LastBuild == false && heart.Count == 0 {
		heart.Count = 3
	} else if message.result.pass == false {
		heart.Count--
		if heart.Count < 0 {
			heart.Count = 0
		}
	}

	heart.LastBuild = message.result.pass
	heart.LastBuilderID = user.ID

	// user

	if message.result.pass == true {
		user.Exp++
		user.Streak++
	} else {
		user.Streak = 0
	}

	if heart.Count == 0 {
		user.Deaths++
	}

	user.CalculateLevel()
	updateBadges(user, kill, save)

}
