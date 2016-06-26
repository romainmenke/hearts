package fakedb

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/net/context"
	"limbo.services/trace"
)

func (db *FakeDB) Persist(ctx context.Context) error {

	span, ctx := trace.New(ctx, "Persist Database")
	defer span.Close()

	err := add()
	if err != nil {
		return span.Error(err)
	}
	err = commit()
	if err != nil {
		return span.Error(err)
	}
	err = push()
	if err != nil {
		return span.Error(err)
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

func push() error {

	user := os.Getenv("GIT_USER")
	password := os.Getenv("GIT_PASS")
	url := fmt.Sprintf("https://%s:%s@github.com/romainmenke/hearts.git", user, password)

	cmd := exec.Command("git", "push", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
