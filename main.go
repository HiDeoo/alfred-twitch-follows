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
		items, err = getFollowedStreamItems(twitch.GetFollowedStreams)
	} else {
		items, err = getFollowItems(twitch.GetFollows)
	}

	if err != nil {
		alfred.SendError(err)

		return
	}

	alfred.SendResult(items)
}

func getFollowItems(getter func() ([]twitch.Follow, error)) ([]alfred.Item, error) {
	follows, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(follows))

	for i, follow := range follows {
		url := fmt.Sprintf("https://www.twitch.tv/%s", follow.ToLogin)

		items[i] = alfred.Item{
			Title:    follow.ToName,
			SubTitle: url,
			Arg:      url,
		}
	}

	return items, nil
}

func getFollowedStreamItems(getter func() ([]twitch.Stream, error)) ([]alfred.Item, error) {
	streams, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(streams))

	for i, stream := range streams {
		items[i] = alfred.Item{
			Title:    stream.UserName,
			SubTitle: fmt.Sprintf("%s - %d viewers - %s", stream.GameName, stream.ViewerCount, stream.Title),
			Arg:      fmt.Sprintf("https://www.twitch.tv/%s", stream.UserLogin),
		}
	}

	return items, nil
}
