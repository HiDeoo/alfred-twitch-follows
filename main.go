package main

import (
	"flag"
	"fmt"

	"github.com/HiDeoo/alfred-twitch-follows/alfred"
	"github.com/HiDeoo/alfred-twitch-follows/twitch"
)

func main() {
	returnLiveFollows := flag.Bool("live", false, "return only live follows")
	flag.Parse()

	var items []alfred.Item
	var err error

	if *returnLiveFollows {
		items, err = getFollowedStreamItems()
	} else {
		items, err = getFollowItems()
	}

	if err != nil {
		alfred.SendError(err)

		return
	}

	alfred.SendResult(items)
}

func getFollowItems() ([]alfred.Item, error) {
	follows, err := twitch.GetFollows()

	if err != nil {
		return nil, err
	}

	return mapFollowsToItems(follows), nil
}

func getFollowedStreamItems() ([]alfred.Item, error) {
	streams, err := twitch.GetFollowedStreams()

	if err != nil {
		return nil, err
	}

	return mapStreamsToItems(streams), nil
}

func mapFollowsToItems(follows []twitch.Follow) []alfred.Item {
	items := make([]alfred.Item, len(follows))

	for i, follow := range follows {
		url := fmt.Sprintf("https://www.twitch.tv/%s", follow.ToLogin)

		items[i] = alfred.Item{
			Title:    follow.ToName,
			SubTitle: url,
			Arg:      url,
		}
	}

	return items
}

func mapStreamsToItems(streams []twitch.Stream) []alfred.Item {
	items := make([]alfred.Item, len(streams))

	for i, stream := range streams {
		items[i] = alfred.Item{
			Title:    stream.UserName,
			SubTitle: fmt.Sprintf("%s - %d viewers - %s", stream.GameName, stream.ViewerCount, stream.Title),
			Arg:      fmt.Sprintf("https://www.twitch.tv/%s", stream.UserLogin),
		}
	}

	return items
}
