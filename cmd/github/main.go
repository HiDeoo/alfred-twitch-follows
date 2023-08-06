package main

import (
	"fmt"
	"time"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/alfred"
	timeago "github.com/caarlos0/timea.go"
)

func main() {
	var items []alfred.Item
	var err error

	items, err = getRepos(GetAllRepos)

	if err != nil {
		alfred.SendError(err)

		return
	}

	alfred.SendResult(items, alfred.Item{
		BaseItem: alfred.BaseItem{Title: "You have no repositories! ¯\\_(ツ)_/¯", SubTitle: "Start coding…"},
		Arg:      "https://github.com/new",
	})
}

func getRepos(getter func() ([]GHRepo, error)) ([]alfred.Item, error) {
	repos, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(repos))

	for i, repo := range repos {
		pushedAt, err := time.Parse(time.RFC3339, repo.PushedAt)

		if err != nil {
			return nil, err
		}

		items[i] = alfred.Item{
			BaseItem: alfred.BaseItem{
				Title:    repo.FullName,
				SubTitle: fmt.Sprintf("Last activity %s", timeago.Of(pushedAt)),
			},
			Arg: repo.HtmlURL,
		}
	}

	return items, nil
}
