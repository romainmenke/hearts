package main

import (
	"fmt"

	"golang.org/x/net/context"

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

	fmt.Println("server.loadingDB")

	db = fakedb.New("/go/src/github.com/romainmenke/hearts/db/", "/go/src/github.com/romainmenke/hearts/db/")
	db.LoadGit(context.Background())

	trace.DefaultHandler = dev.NewHandler(nil)

	go serveHTTP()

	serveRCP()

}
