package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/romainmenke/hearts/pkg/fakedb"
	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"limbo.services/trace"
)

func serveRCP() {

	fmt.Println("server.grpc.starting")

	fmt.Println("server.grpc.loadingDB")

	db = fakedb.New("/go/src/app/db/", "/go/src/app/db/")
	db.LoadGit(context.Background())

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	fmt.Println("server.grpc.listening")

	s := grpc.NewServer()

	wercker.RegisterNotificationServiceServer(s, &server{})

	fmt.Println("server.grpc.ready")

	s.Serve(lis)

}

type server struct{}

func (s *server) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {

	if in == nil || in.Git == nil || in.Build == nil {
		return &wercker.WerckerResponse{Success: false}, errors.New("nil message")
	}

	if in.Git.Branch != "master" {
		return &wercker.WerckerResponse{Success: true}, nil
	}

	err := s.update(ctx, in)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	return &wercker.WerckerResponse{Success: true}, nil
}

func (s *server) update(ctx context.Context, message *wercker.WerckerMessage) error {

	span, ctx := trace.New(ctx, "server.grpc.handle.wercker")
	defer span.Close()

	heart := s.loadHeart(ctx, message)
	user := s.loadUser(ctx, message)

	applyChanges(ctx, message, heart, user)

	err := s.saveHeart(ctx, message, heart)
	if err != nil {
		return span.Error(err)
	}

	err = s.saveUser(ctx, message, user)
	if err != nil {
		return span.Error(err)
	}

	err = db.Persist(ctx)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func (s *server) loadHeart(ctx context.Context, message *wercker.WerckerMessage) *fakedb.Heart {

	span, ctx := trace.New(ctx, "server.grpc.loadHeart")
	defer span.Close()

	heart, err := db.LoadHeart(ctx, message.Git.Domain, message.Git.Owner, message.Git.Repository)
	if err != nil || heart == nil {
		span.Error(err)

		pass := message.Result.Result
		newHeart := &fakedb.Heart{
			Count:     3,
			LastBuild: pass,
			Domain:    message.Git.Domain,
			Owner:     message.Git.Owner,
			Repo:      message.Git.Repository,
		}
		return newHeart
	}

	return heart
}

func (s *server) loadUser(ctx context.Context, message *wercker.WerckerMessage) *fakedb.User {

	span, ctx := trace.New(ctx, "server.grpc.loadUser")
	defer span.Close()

	user, err := db.LoadUser(ctx, message.Git.Domain, message.Build.User)
	if err != nil || user == nil {
		span.Error(err)

		newUser := &fakedb.User{
			Domain: message.Git.Domain,
			Name:   message.Build.User,
			Level:  0,
			Exp:    0,
			Streak: 0,
			Deaths: 0,
		}
		return newUser
	}

	return user
}

func applyChanges(ctx context.Context, message *wercker.WerckerMessage, heart *fakedb.Heart, user *fakedb.User) {

	span, ctx := trace.New(ctx, "server.grpc.applyChanges")
	defer span.Close()

	// HEART
	if message.Result.Result == true && heart.LastBuild == false && heart.Count != 0 {
		heart.Count++
		if heart.Count > 3 {
			heart.Count = 3
		}
	} else if message.Result.Result == true && heart.LastBuild == false && heart.Count == 0 {
		heart.Count = 3
	} else if message.Result.Result == false {
		heart.Count--
		if heart.Count < 0 {
			heart.Count = 0
		}
	}

	heart.LastBuild = message.Result.Result
	heart.LastBuilder = user

	// user

	if message.Result.Result == true {
		user.Exp++
		user.Streak++
	} else {
		user.Streak = 0
	}

	if heart.Count == 0 {
		user.Deaths++
	}

}

func (s *server) saveHeart(ctx context.Context, message *wercker.WerckerMessage, heart *fakedb.Heart) error {

	span, ctx := trace.New(ctx, "server.grpc.saveHeart")
	defer span.Close()

	err := db.SaveObject(ctx, heart)
	if err != nil {
		return span.Error(err)
	}
	return nil

}

func (s *server) saveUser(ctx context.Context, message *wercker.WerckerMessage, user *fakedb.User) error {

	span, ctx := trace.New(ctx, "server.grpc.saveUser")
	defer span.Close()

	err := db.SaveObject(ctx, user)
	if err != nil {
		return span.Error(err)
	}
	return nil

}
