package main

import (
	"fmt"

	"github.com/romainmenke/hearts/pkg/fakedb"
	"limbo.services/trace"
	"limbo.services/trace/dev"
)

const (
	port = ":50051"
)

var (
	db *fakedb.FakeDB
)

func main() {

	fmt.Println("server.starting")

	trace.DefaultHandler = dev.NewHandler(nil)

	go serveHTTP()

	serveRCP()

}
