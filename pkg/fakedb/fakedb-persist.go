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

	span, _ := trace.New(ctx, "Persist Database")
	defer span.Close()

	err := cd(db.gitRoot)
	if err != nil {
		return span.Error(err)
	}

	err = add()
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

func cd(path string) error {

	cmd := exec.Command("cd", path)
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

	cmd := exec.Command("git", "commit", "-m", "'entry'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func push() error {

	user := os.Getenv("USER")
	password := os.Getenv("PASS")
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

func setupUser() error {
	user := "heartsbot"

	user = fmt.Sprintf("'%s'", user)

	cmd := exec.Command("git", "config", "user.name", user)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil

}

func setupEmail() error {
	email := "romainmenke+heartsbot@gmail.com"

	email = fmt.Sprintf("'%s'", email)

	cmd := exec.Command("git", "config", "user.email", email)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
