package main

import (
	"net/http"
	"time"

	"github.com/HiDeoo/alfred-twitch-follows/alfred"
	"github.com/HiDeoo/alfred-twitch-follows/twitch"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	follows, err := twitch.GetFollows(httpClient)

	if err != nil {
		alfred.SendError(err)
	}

	alfred.SendResult(MapFollowsToItems(follows))
}

func MapFollowsToItems(from []twitch.Follow) []alfred.Item {
	items := make([]alfred.Item, len(from))

	for i, follow := range from {
		items[i] = alfred.Item{
			Title:    follow.ToName,
			SubTitle: "test",
		}
	}

	return items
}
