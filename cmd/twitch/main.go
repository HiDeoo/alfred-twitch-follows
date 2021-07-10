package main

import (
	"flag"
	"fmt"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/alfred"
)

func main() {
	returnLiveFollows := flag.Bool("live", false, "return only live follows")
	flag.Parse()

	var items []alfred.Item
	var err error

	if *returnLiveFollows {
		items, err = getFollowedStreamItems(GetFollowedStreams)
	} else {
		items, err = getFollowItems(GetFollows)
	}

	if err != nil {
		alfred.SendError(err)

		return
	}

	alfred.SendResult(items, alfred.Item{
		BaseItem: alfred.BaseItem{Title: "You're alone! ¯\\_(ツ)_/¯", SubTitle: "Try browsing Twitch…"},
		Arg:      "https://www.twitch.tv/directory/following",
	})
}

func getFollowItems(getter func() ([]TwitchFollow, error)) ([]alfred.Item, error) {
	follows, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(follows))

	for i, follow := range follows {
		url := fmt.Sprintf("https://www.twitch.tv/%s", follow.ToLogin)

		items[i] = alfred.Item{
			BaseItem: alfred.BaseItem{
				Title:    follow.ToName,
				SubTitle: url,
			},
			Arg: url,
		}
	}

	return items, nil
}

func getFollowedStreamItems(getter func() ([]TwitchStream, error)) ([]alfred.Item, error) {
	streams, err := getter()

	if err != nil {
		return nil, err
	}

	items := make([]alfred.Item, len(streams))

	for i, stream := range streams {
		items[i] = alfred.Item{
			BaseItem: alfred.BaseItem{
				Title:    stream.UserName,
				SubTitle: fmt.Sprintf("%s - %d viewers - %s", stream.GameName, stream.ViewerCount, stream.Title),
			},
			Arg: fmt.Sprintf("https://www.twitch.tv/%s", stream.UserLogin),
		}
	}

	return items, nil
}
