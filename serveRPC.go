package main

import (
	"fmt"
	"log"
	"net"

	"limbo.services/trace"

	"github.com/romainmenke/universal-notifier/pkg/wercker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func serveRCP() {

	fmt.Println("server.tcp.listening on port : " + port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	fmt.Println("server.grpc.startingServer")

	s := grpc.NewServer()

	wercker.RegisterNotificationServiceServer(s, &grpcServer{})

	fmt.Println("server.grpc.ready")

	s.Serve(lis)

}

type grpcServer struct{}

func (s *grpcServer) Notify(ctx context.Context, in *wercker.WerckerMessage) (*wercker.WerckerResponse, error) {

	span, _ := trace.New(ctx, "server.grpc.notify")
	defer span.Close()

	if in == nil || in.Git == nil || in.Build == nil {
		err := span.Error("nil message")
		return &wercker.WerckerResponse{Success: false}, err
	}

	if in.Git.Branch != "master" {
		span.Log("not on the main branch")
		return &wercker.WerckerResponse{Success: true}, nil
	}

	message := newFromWercker(in)

	err := update(ctx, db, message)
	if err != nil {
		span.Error(err)
		return &wercker.WerckerResponse{Success: false}, err
	}

	return &wercker.WerckerResponse{Success: true}, nil
}
