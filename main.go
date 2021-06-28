package main

import (
	"fmt"

	"github.com/HiDeoo/alfred-twitch-follows/alfred"
	"github.com/HiDeoo/alfred-twitch-follows/twitch"
)

func main() {
	follows, err := twitch.GetFollows()

	if err != nil {
		alfred.SendError(err)
	}

	alfred.SendResult(mapFollowsToItems(follows))
}

func mapFollowsToItems(from []twitch.Follow) []alfred.Item {
	items := make([]alfred.Item, len(from))

	for i, follow := range from {
		url := fmt.Sprintf("https://www.twitch.tv/%s", follow.ToLogin)

		items[i] = alfred.Item{
			Title:    follow.ToName,
			SubTitle: url,
			Arg:      url,
		}
	}

	return items
}
