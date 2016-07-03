package fakedb

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
)

type FakeSheme interface {
	Path() string
	Filename() string
	Bytes() (*[]byte, error)
	Hash() int
}

type Heart struct {
	ID            string
	Count         int
	LastBuild     bool
	LastBuilderID string
	Owner         string
	Domain        string
	Repo          string
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

func (h *Heart) Bytes() (*[]byte, error) {
	b, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (h *Heart) Hash() int {
	comp := h.ID + strconv.Itoa(h.Count) + strconv.FormatBool(h.LastBuild) + h.LastBuilderID + h.Owner + h.Domain + h.Repo
	hash := fnv.New64()
	hash.Write([]byte(comp))
	return int(hash.Sum64())
}

func (h *Heart) FullPath() string {
	return fmt.Sprintf("heart/%s/%s/%s", h.Domain, h.Owner, h.Repo)
}

func (h *Heart) SVGFileName() string {
	r := fmt.Sprintf("%s.svg", h.Repo)
	return r
}

func (h *Heart) SVG() string {
	return svg(h.Count)
}

type User struct {
	ID     string
	Domain string
	Name   string
	Level  int
	Exp    int
	Streak int
	Deaths int
	Badges []*Badge
}

func (u *User) CalculateLevel() {

	f := float64(u.Exp)
	f = f * 0.2345
	s := math.Sqrt(f)
	u.Level = int(s)

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

func (u *User) Hash() int {
	comp := u.ID + u.Domain + u.Name + strconv.Itoa(u.Level) + strconv.Itoa(u.Exp) + strconv.Itoa(u.Streak) + strconv.Itoa(u.Deaths)
	for _, badge := range u.Badges {
		comp += badge.HashString()
	}
	hash := fnv.New64()
	hash.Write([]byte(comp))
	return int(hash.Sum64())
}

func (u *User) FullPath() string {
	return fmt.Sprintf("user/%s/%s", u.Domain, u.Name)
}

type Badge struct {
	Class    int
	Name     string
	Progress int
}

func (b *Badge) HashString() string {
	return strconv.Itoa(b.Class) + strconv.Itoa(b.Progress) + b.Name
}
