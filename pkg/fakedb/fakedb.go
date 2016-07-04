package fakedb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
	"limbo.services/trace"
)

type FakeDB struct {
	dbRoot  string
	gitRoot string
	gitUser string
	gitPass string
}

func New(dbRoot string, gitRoot string, gitUser string, gitPass string) *FakeDB {
	db := FakeDB{dbRoot: dbRoot, gitRoot: gitRoot, gitUser: gitUser, gitPass: gitPass}
	return &db
}

func (db *FakeDB) SaveObject(ctx context.Context, object FakeSheme) error {

	span, ctx := trace.New(ctx, "fakedb.saveObject")
	defer span.Close()

	b, err := object.Bytes()
	if err != nil {
		return span.Error(err)
	}

	err = db.Save(ctx, b, object.Path(), object.Filename())
	if err != nil {
		return span.Error(err)
	}
	return nil
}

func (db *FakeDB) SaveSVG(ctx context.Context, heart Heart) error {

	span, ctx := trace.New(ctx, "fakedb.saveSVG")
	defer span.Close()

	svgString := svg(heart.Count)
	b := []byte(svgString)

	err := db.Save(ctx, &b, heart.Path(), heart.SVGFileName())
	if err != nil {
		return span.Error(err)
	}
	return nil
}

func (db *FakeDB) Save(ctx context.Context, data *[]byte, dir string, filename string) error {

	span, ctx := trace.New(ctx, "fakedb.saveBytes")
	defer span.Close()

	dirPath := fmt.Sprint(db.dbRoot, dir)
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dirPath, 0755)
		} else {
			return span.Error(err)
		}
	}

	path := fmt.Sprint(db.dbRoot, dir, filename)
	os.Remove(path)

	var file, err = os.Create(path)
	if err != nil {
		return span.Error(err)
	}
	defer file.Close()

	err = ioutil.WriteFile(path, *data, 0644)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func (db *FakeDB) LoadHeart(ctx context.Context, domain string, owner string, repo string) (*Heart, error) {

	span, ctx := trace.New(ctx, "fakedb.loadHeart")
	defer span.Close()

	heart := Heart{Domain: domain, Owner: owner, Repo: repo}
	b, err := db.Load(ctx, heart.Path(), heart.Filename())
	if err != nil || b == nil {
		return nil, span.Error(err)
	}

	err = json.Unmarshal(*b, &heart)
	if err != nil {
		return nil, span.Error(err)
	}

	return &heart, nil
}

func (db *FakeDB) LoadUser(ctx context.Context, domain string, name string) (*User, error) {

	span, ctx := trace.New(ctx, "fakedb.loadUser")
	defer span.Close()

	user := User{Domain: domain, Name: name}
	b, err := db.Load(ctx, user.Path(), user.Filename())
	if err != nil || b == nil {
		return nil, span.Error(err)
	}

	err = json.Unmarshal(*b, &user)
	if err != nil {
		return nil, span.Error(err)
	}

	return &user, nil
}

func (db *FakeDB) Load(ctx context.Context, dir string, filename string) (*[]byte, error) {

	span, ctx := trace.New(ctx, "fakedb.loadObject")
	defer span.Close()

	path := fmt.Sprint(db.dbRoot, dir, filename)

	if _, err := os.Stat(path); err != nil {
		return nil, span.Error(err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, span.Error(err)
	}

	return &b, nil
}
