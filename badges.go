package main

import "github.com/romainmenke/hearts/pkg/fakedb"

func immortalBadge(u *fakedb.User, badge *fakedb.Badge, kill bool) {
	badge.Name = "Immortal"
	badge.Class = 3
	if kill {
		badge.Progress = 0
	} else {
		badge.Progress++
	}
}

func saviourBadge(u *fakedb.User, badge *fakedb.Badge, save bool) {
	badge.Name = "Saviour"
	badge.Class = 2
	if save {
		badge.Progress++
	}
}

func destoyerBadge(u *fakedb.User, badge *fakedb.Badge, kill bool) {
	badge.Name = "Destroyer"
	badge.Class = 2
	if kill {
		badge.Progress++
	}
}

func updateBadges(u *fakedb.User, kill bool, save bool) {

	immortal := fakedb.Badge{
		Name:  "Immortal",
		Class: 2,
	}

	saviour := fakedb.Badge{
		Name:  "Saviour",
		Class: 2,
	}

	destoyer := fakedb.Badge{
		Name:  "Destoyer",
		Class: 2,
	}

	for _, badge := range u.Badges {
		switch badge.Name {
		case immortal.Name:
			immortal.Progress = badge.Progress
			immortalBadge(u, &immortal, kill)
		case destoyer.Name:
			destoyer.Progress = badge.Progress
			destoyerBadge(u, &destoyer, kill)
		case saviour.Name:
			saviour.Progress = badge.Progress
			saviourBadge(u, &saviour, save)
		}
	}

	u.Badges = []fakedb.Badge{immortal, saviour, destoyer}

}
