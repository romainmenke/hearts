package fakedb

import (
	"encoding/json"
	"fmt"
)

type FakeSheme interface {
	Path() string
	Filename() string
	Bytes() (*[]byte, error)
}

type Heart struct {
	Count     int
	LastBuild bool
	Owner     string
	Domain    string
	Repo      string
}

func (h *Heart) Path() string {
	d := fmt.Sprintf("heart/%s/", h.Domain)
	o := fmt.Sprintf("%s/", h.Owner)
	return fmt.Sprint(d, o)
}

func (h *Heart) Filename() string {
	r := fmt.Sprintf("%s.json", h.Repo)
	return r
}

func (h *Heart) SVGFileName() string {
	r := fmt.Sprintf("%s.svg", h.Repo)
	return r
}

func (h *Heart) Bytes() (*[]byte, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

type User struct {
	Domain string
	Name   string
	Level  int
	Exp    int
	Streak int
	Deaths int
}

func (u *User) Path() string {
	d := fmt.Sprintf("user/%s/", u.Domain)
	return d
}

func (u *User) Filename() string {
	r := fmt.Sprintf("%s.json", u.Name)
	return r
}

func (u *User) Bytes() (*[]byte, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
