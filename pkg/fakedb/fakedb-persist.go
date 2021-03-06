package fakedb

import (
	"bytes"
	"fmt"
	"os/exec"

	"golang.org/x/net/context"
	"limbo.services/trace"
)

func (db *FakeDB) Persist(ctx context.Context) error {

	span, ctx := trace.New(ctx, "fakedb.persist")
	defer span.Close()

	err := add()
	if err != nil {
		return nil
	}
	err = commit()
	if err != nil {
		return nil
	}
	err = push(db.gitUser, db.gitPass)
	if err != nil {
		return span.Error(err)
	}

	return nil
}

func (db *FakeDB) LoadGit(ctx context.Context) error {

	span, ctx := trace.New(ctx, "fakedb.load")
	defer span.Close()

	err := pull()
	if err != nil {
		return nil
	}

	return nil

}

func pull() error {

	cmd := exec.Command("git", "pull")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func add() error {

	cmd := exec.Command("git", "add", ".")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func commit() error {

	cmd := exec.Command("git", "commit", "-m", "entry")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func push(user string, pass string) error {

	url := fmt.Sprintf("https://%s:%s@github.com/romainmenke/hearts.git", user, pass)

	cmd := exec.Command("git", "push", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
