package main

import (
	"github.com/HiDeoo/alfred-workflow-tools/pkg/alfred"
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
	var repos []GHRepo
	var err error
	var updateCache = false

	alfred.GetCache(&repos)

	if repos == nil {
		updateCache = true
		repos, err = getter()
	}

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(repos))

	for i, repo := range repos {
		if err != nil {
			return nil, err
		}

		items[i] = alfred.Item{
			BaseItem: alfred.BaseItem{
				Title:    repo.FullName,
				SubTitle: repo.HtmlURL,
			},
			Arg: repo.HtmlURL,
		}
	}

	if updateCache {
		alfred.SetCache(repos)
	}

	return items, nil
}
