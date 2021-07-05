package main

import (
	"fmt"
	"strconv"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/alfred"
)

func main() {
	items, err := getUnwatchedShowItems(getCurrentShowsWithUnwatchedEpisodes)

	if err != nil {
		alfred.SendError(err)

		return
	}

	alfred.SendResult(items)
}

func getUnwatchedShowItems(getter func() ([]BSShow, error)) ([]alfred.Item, error) {
	shows, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(shows))

	for i, show := range shows {
		items[i] = alfred.Item{
			BaseItem: alfred.BaseItem{
				Title:    show.Title,
				SubTitle: fmt.Sprintf("%d episodes remaining (%s total)", show.User.Remaining, show.Episodes),
			},
			Arg: strconv.Itoa(show.ID),
		}
	}

	return items, nil
}
