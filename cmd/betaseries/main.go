package main

import (
	"flag"

	"fmt"
	"strconv"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/alfred"
)

func main() {
	showID := flag.String("watched", "", "mark all episodes of a show as watched")
	flag.Parse()

	shouldMarkShowAsWatched := *showID != ""

	var items []alfred.Item
	var err error

	if shouldMarkShowAsWatched {
		err = markShowAsWatched(*showID)
	} else {
		items, err = getUnwatchedShowItems(getCurrentShowsWithUnwatchedEpisodes)
	}

	if err != nil {
		alfred.SendError(err)

		return
	}

	if !shouldMarkShowAsWatched {
		alfred.SendResult(items)
	}
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
