package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/request"
)

var client request.Client

func init() {
	client = request.NewClient("https://api.twitch.tv/helix/")
	client.SetHeaders(http.Header{
		"Client-ID":     []string{os.Getenv("TWITCH_CLIENT_ID")},
		"Authorization": []string{"Bearer " + os.Getenv("TWITCH_OAUTH_TOKEN")},
	})
}

func GetFollows() ([]TwitchFollow, error) {
	currentUser, err := getCurrentUser()

	if err != nil {
		return nil, err
	}

	var cursor = ""
	var allFollows []TwitchFollow

	for {
		follows, err := getFollowsWithPagination(currentUser.Id, cursor)

		if err != nil {
			return nil, err
		}

		allFollows = append(allFollows, follows.Data...)

		if len(follows.Pagination.Cursor) == 0 {
			break
		}

		cursor = follows.Pagination.Cursor
	}

	return allFollows, nil
}

func GetFollowedStreams() ([]TwitchStream, error) {
	currentUser, err := getCurrentUser()

	if err != nil {
		return nil, err
	}

	var cursor = ""
	var allStreams []TwitchStream

	for {
		streams, err := getFollowedStreamsWithPagination(currentUser.Id, cursor)

		if err != nil {
			return nil, err
		}

		allStreams = append(allStreams, streams.Data...)

		if len(streams.Pagination.Cursor) == 0 {
			break
		}

		cursor = streams.Pagination.Cursor
	}

	return allStreams, nil
}

func GetGameStreams(game string, lang string) ([]TwitchStream, error) {
	var cursor = ""
	var allStreams []TwitchStream

	for {
		streams, err := getGameStreamsWithPagination(game, lang, cursor)

		if err != nil {
			return nil, err
		}

		allStreams = append(allStreams, streams.Data...)

		if len(streams.Pagination.Cursor) == 0 {
			break
		}

		cursor = streams.Pagination.Cursor
	}

	return allStreams, nil
}

func getFollowsWithPagination(userID string, cursor string) (*TwitchFollows, error) {
	queryParams := url.Values{}
	queryParams.Set("from_id", userID)
	queryParams.Set("first", "100")
	queryParams.Set("after", cursor)

	res, err := query(client.Get("users/follows", queryParams))

	if err != nil {
		return nil, err
	}

	follows := TwitchFollows{}

	if err = json.Unmarshal(res.Data, &follows); err != nil {
		return nil, err
	}

	return &follows, nil
}

func getFollowedStreamsWithPagination(userID string, cursor string) (*TwitchStreams, error) {
	queryParams := url.Values{}
	queryParams.Set("user_id", userID)
	queryParams.Set("first", "100")
	queryParams.Set("after", cursor)

	res, err := query(client.Get("streams/followed", queryParams))

	if err != nil {
		return nil, err
	}

	streams := TwitchStreams{}

	if err = json.Unmarshal(res.Data, &streams); err != nil {
		return nil, err
	}

	return &streams, nil
}

func getGameStreamsWithPagination(game string, lang string, cursor string) (*TwitchStreams, error) {
	queryParams := url.Values{}
	queryParams.Set("first", "100")
	queryParams.Set("type", "live")
	queryParams.Set("game_id", game)
	queryParams.Set("after", cursor)

	if len(lang) > 0 {
		queryParams.Set("language", lang)
	}

	res, err := query(client.Get("streams", queryParams))

	if err != nil {
		return nil, err
	}

	streams := TwitchStreams{}

	if err = json.Unmarshal(res.Data, &streams); err != nil {
		return nil, err
	}

	return &streams, nil
}

func getCurrentUser() (*TwitchUser, error) {
	res, err := query(client.Get("users", nil))

	if err != nil {
		return nil, err
	}

	users := TwitchUsers{}

	if err = json.Unmarshal(res.Data, &users); err != nil {
		return nil, err
	}

	if len(users.Data) != 1 {
		return nil, errors.New("unable to get current user")
	}

	return &users.Data[0], nil
}

func query(res *request.Response, err error) (*request.Response, error) {
	if res.StatusCode != 200 {
		twitchError := &TwitchError{}

		if err := json.Unmarshal(res.Data, &twitchError); err != nil || twitchError.Error == "" {
			return nil, fmt.Errorf("unable to fetch Twitch data (error: %d)", res.StatusCode)
		}

		return nil, errors.New(twitchError.Message)
	}

	return res, err
}
