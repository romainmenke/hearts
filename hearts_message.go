package main

import "github.com/romainmenke/universal-notifier/pkg/wercker"
import "github.com/romainmenke/travis"

type incomingMessage struct {
	repo   repo
	user   user
	result result
}

type repo struct {
	owner  string
	domain string
	name   string
}

type user struct {
	name string
}

type result struct {
	pass bool
}

func newFromWercker(w *wercker.WerckerMessage) *incomingMessage {

	return &incomingMessage{
		repo: repo{
			owner:  w.Git.Owner,
			domain: w.Git.Domain,
			name:   w.Git.Repository,
		},
		result: result{
			pass: w.Result.Result,
		},
		user: user{
			name: w.Build.User,
		},
	}
}

func newFromTravis(t *travis.PayloadObject) *incomingMessage {

	return &incomingMessage{
		repo: repo{
			owner:  t.Repository.OwnerName,
			domain: t.Repository.Domain(),
			name:   t.Repository.Name,
		},
		result: result{
			pass: t.StatusMessage == "Passed",
		},
		user: user{
			name: t.AuthorName,
		},
	}
}
