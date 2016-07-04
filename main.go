package main

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/romainmenke/hearts/pkg/fakedb"
	"github.com/romainmenke/hearts/pkg/memcache"
	"limbo.services/trace"
	"limbo.services/trace/dev"
)

const (
	port = ":50051"
)

var (
	cache *memcache.MemCache
)

func main() {

	trace.DefaultHandler = dev.NewHandler(nil)

	fmt.Println("server.starting")
	fmt.Println("server.loadingDB")

	gitUser := os.Getenv("GIT_USER")
	gitPass := os.Getenv("GIT_PASS")
	if gitUser == "" || gitPass == "" {
		fmt.Println("no git credentials")
	}
	db := fakedb.New("/go/src/github.com/romainmenke/hearts/db/", "/go/src/github.com/romainmenke/hearts/db/", gitUser, gitPass)
	db.LoadGit(context.Background())

	cache = memcache.New(db)

	memcache.RunCacheWorker(cache)
	memcache.RunPersistWorker(cache)

	go serveHTTP()

	serveRCP()

}
