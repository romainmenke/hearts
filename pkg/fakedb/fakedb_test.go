package fakedb

import (
	"fmt"
	"testing"

	"limbo.services/trace"
	"limbo.services/trace/dev"

	"golang.org/x/net/context"
)

func TestWriteHeart(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	heart := Heart{
		Count:     2,
		LastBuild: true,
		Owner:     "romainmenke",
		Domain:    "github.com",
		Repo:      "hearts",
	}

	err := db.SaveObject(ctx, &heart)
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}
}

func TestWriteSVG(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	heart := Heart{
		Count:     2,
		LastBuild: true,
		Owner:     "romainmenke",
		Domain:    "github.com",
		Repo:      "hearts",
	}

	err := db.SaveSVG(ctx, heart)
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}
}

func TestReadHeart(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	loadHeart, err := db.LoadHeart(ctx, "github.com", "romainmenke", "hearts")
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}

	compareHeart := Heart{
		Count:     2,
		LastBuild: true,
		Owner:     "romainmenke",
		Domain:    "github.com",
		Repo:      "hearts",
	}

	if *loadHeart != compareHeart {
		fmt.Print(*loadHeart)
		t.Fail()
	}
}

func TestReadNonHeart(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	loadHeart, err := db.LoadHeart(ctx, "github.com", "blah", "hearts")

	if loadHeart != nil {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}

func TestWriteUser(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	user := User{
		Domain: "github.com",
		Name:   "romainmenke",
		Level:  10,
		Exp:    0,
		Streak: 0,
		Deaths: 999,
	}

	err := db.SaveObject(ctx, &user)
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}
}

func TestReadUser(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	loadUser, err := db.LoadUser(ctx, "github.com", "romainmenke")
	if err != nil {
		fmt.Print(err)
		t.Fail()
	}

	compareUser := User{
		Domain: "github.com",
		Name:   "romainmenke",
		Level:  10,
		Exp:    0,
		Streak: 0,
		Deaths: 999,
	}

	if loadUser.Hash() != compareUser.Hash() {
		fmt.Print(*loadUser)
		t.Fail()
	}
}

func TestReadNonUser(t *testing.T) {

	trace.DefaultHandler = dev.NewHandler(nil)
	ctx := context.Background()

	db := New("/Users/romainmenke/Go/src/github.com/romainmenke/hearts/pkg/fakedb/testdb/", "/Users/romainmenke/Go/src/github.com/romainmenke/hearts/", "heartsbot", "nopass")

	loadUser, err := db.LoadUser(ctx, "github.com", "blah")

	if loadUser != nil {
		t.Fail()
	}

	if err == nil {
		t.Fail()
	}
}
