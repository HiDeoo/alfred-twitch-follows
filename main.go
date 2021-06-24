package main

import (
	"fmt"
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

	fmt.Println(follows)
}
