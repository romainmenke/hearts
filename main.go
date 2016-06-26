package main

import (
	"fmt"
	"log"
	"net"

	"limbo.services/trace"
	"limbo.services/trace/dev"

	"golang.org/x/net/context"

	"github.com/romainmenke/hearts/pkg/fakedb"
	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var (
	db = fakedb.New("/go/src/app/", "/go/src/app/")
)

func main() {

	trace.DefaultHandler = dev.NewHandler(nil)

	fmt.Println("Starting Hearts")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	fmt.Printf("Listening on port : %s", port)

	s := grpc.NewServer()
	srv := server{}

	fmt.Println("Loading DB Connection")

	wercker.RegisterNotificationServiceServer(s, &srv)

	fmt.Println("Ready to serve clients")

	s.Serve(lis)

}

type server struct{}

func (s *server) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {

	heart, err := s.heart(ctx, in)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = db.SaveSVG(ctx, *heart)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = s.user(ctx, in, heart)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	err = db.Persist(ctx)
	if err != nil {
		return &wercker.WerckerResponse{Success: false}, err
	}

	return &wercker.WerckerResponse{Success: true}, nil
}

func (s *server) heart(ctx context.Context, message *wercker.WerckerMessage) (*fakedb.Heart, error) {

	span, ctx := trace.New(ctx, "Update Heart")
	defer span.Close()

	var heart *fakedb.Heart

	heart, err := db.LoadHeart(ctx, message.Git.Domain, message.Git.Owner, message.Git.Repository)
	if err != nil || heart == nil {
		span.Error(err)

		newH, newErr := s.newHeart(ctx, message)
		if newErr != nil {
			return nil, span.Error(newErr)
		}
		return newH, nil
	}

	if message.Result.Result == true && heart.LastBuild == false {
		heart.Count++
		if heart.Count > 3 {
			heart.Count = 3
		}
	} else if message.Result.Result == false {
		heart.Count--
		if heart.Count < 0 {
			heart.Count = 0
		}
	}

	heart.LastBuild = message.Result.Result

	err = db.SaveObject(ctx, heart)
	if err != nil {
		return nil, span.Error(err)
	}

	return heart, nil
}

func (s *server) newHeart(ctx context.Context, message *wercker.WerckerMessage) (*fakedb.Heart, error) {

	span, ctx := trace.New(ctx, "New Heart")
	defer span.Close()

	pass := message.Result.Result
	heart := &fakedb.Heart{
		Count:     3,
		LastBuild: pass,
		Domain:    message.Git.Domain,
		Owner:     message.Git.Owner,
		Repo:      message.Git.Repository,
	}

	err := db.SaveObject(ctx, heart)
	if err != nil {
		return nil, span.Error(err)
	}

	return heart, nil
}

func (s *server) user(ctx context.Context, message *wercker.WerckerMessage, heart *fakedb.Heart) error {

	span, ctx := trace.New(ctx, "Update User")
	defer span.Close()

	var user *fakedb.User

	user, err := db.LoadUser(ctx, message.Git.Domain, message.Build.User)
	if err != nil {
		span.Error(err)

		err = s.newUser(ctx, message)
		if err != nil {
			return span.Error(err)
		}
		return nil
	}

	if message.Result.Result == true {
		user.Exp++
		user.Streak++
	} else {
		user.Streak = 0
	}

	if heart.Count == 0 {
		user.Deaths++
	}

	heart.LastBuild = message.Result.Result

	err = db.SaveObject(ctx, heart)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func (s *server) newUser(ctx context.Context, message *wercker.WerckerMessage) error {

	span, ctx := trace.New(ctx, "New User")
	defer span.Close()

	pass := message.Result.Result
	user := &fakedb.User{
		Domain: message.Git.Domain,
		Name:   message.Build.User,
		Level:  0,
		Exp:    0,
		Streak: 0,
		Deaths: 0,
	}

	if pass {
		user.Exp++
		user.Streak++
	}

	err := db.SaveObject(ctx, user)
	if err != nil {
		return span.Error(err)
	}

	return nil
}
