package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFollowItems(t *testing.T) {
	tests := []struct {
		name        string
		followCount int
	}{
		{"ReturnNoFollows", 0},
		{"ReturnSingleFollow", 1},
		{"ReturnMultipleFollows", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var follows = []TwitchFollow{}

			for i := 0; i < test.followCount; i++ {
				followJSON := TwitchFollow{
					FromId:     "123456789",
					FromLogin:  "user",
					ToId:       fmt.Sprintf("1234560%d", i),
					ToLogin:    fmt.Sprintf("user%d", i),
					ToName:     fmt.Sprintf("User%d", i),
					FollowedAt: fmt.Sprintf("2020-03-01T02:23:4%d.009756Z", i),
				}

				follows = append(follows, followJSON)
			}

			items, _ := getFollowItems(func() ([]TwitchFollow, error) {
				return follows, nil
			})

			assert.Equal(t, test.followCount, len(items))

			for i, item := range items {
				assert.Equal(t, follows[i].ToName, item.Title)
				assert.Equal(t, fmt.Sprintf("https://www.twitch.tv/%s", follows[i].ToLogin), item.SubTitle)
				assert.Equal(t, item.SubTitle, item.Arg)
			}
		})
	}
}

func TestGetFollowedStreamItems(t *testing.T) {
	tests := []struct {
		name                 string
		followedStreamsCount int
	}{
		{"ReturnNoFollowedStreams", 0},
		{"ReturnSingleFollowedStream", 1},
		{"ReturnMultipleFollowedStreams", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var streams = []TwitchStream{}

			for i := 0; i < test.followedStreamsCount; i++ {
				streamJSON := TwitchStream{
					Id:           strconv.Itoa(i),
					UserId:       fmt.Sprintf("1234560%d", i),
					UserLogin:    fmt.Sprintf("user%d", i),
					UserName:     fmt.Sprintf("User%d", i),
					GameId:       fmt.Sprintf("game%d", i),
					GameName:     fmt.Sprintf("Game%d", i),
					Type:         "live",
					Title:        fmt.Sprintf("Stream Title %d", i),
					ViewerCount:  i * 100,
					StartedAt:    fmt.Sprintf("2020-02-01T02:23:4%d.009756Z", i),
					Language:     "en",
					ThumbnailUrl: fmt.Sprintf("https://twitch.tv/stream/%d/thumbnail.png", i),
					TagIds:       []string{"1122334455"},
					IsMature:     false,
				}

				streams = append(streams, streamJSON)
			}

			items, _ := getFollowedStreamItems(func() ([]TwitchStream, error) {
				return streams, nil
			})

			assert.Equal(t, test.followedStreamsCount, len(items))

			for i, item := range items {
				assert.Equal(t, streams[i].UserName, item.Title)
				assert.Equal(
					t,
					fmt.Sprintf("%s - %d viewers - %s", streams[i].GameName, streams[i].ViewerCount, streams[i].Title),
					item.SubTitle,
				)
				assert.Equal(t, fmt.Sprintf("https://www.twitch.tv/%s", streams[i].UserLogin), item.Arg)
			}
		})
	}
}

func TestGetGameStreamItems(t *testing.T) {
	tests := []struct {
		name             string
		gameStreamsCount int
	}{
		{"ReturnNoFollowedStreams", 0},
		{"ReturnSingleFollowedStream", 1},
		{"ReturnMultipleFollowedStreams", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var streams = []TwitchStream{}

			for i := 0; i < test.gameStreamsCount; i++ {
				streamJSON := TwitchStream{
					Id:           strconv.Itoa(i),
					UserId:       fmt.Sprintf("1234560%d", i),
					UserLogin:    fmt.Sprintf("user%d", i),
					UserName:     fmt.Sprintf("User%d", i),
					GameId:       fmt.Sprintf("game%d", i),
					GameName:     fmt.Sprintf("Game%d", i),
					Type:         "live",
					Title:        fmt.Sprintf("Stream Title %d", i),
					ViewerCount:  i * 100,
					StartedAt:    fmt.Sprintf("2020-02-01T02:23:4%d.009756Z", i),
					Language:     "en",
					ThumbnailUrl: fmt.Sprintf("https://twitch.tv/stream/%d/thumbnail.png", i),
					TagIds:       []string{"1122334455"},
					IsMature:     false,
				}

				streams = append(streams, streamJSON)
			}

			items, _ := getGameStreamItems(func(game string, lang string) ([]TwitchStream, error) {
				return streams, nil
			}, "game1", "")

			assert.Equal(t, test.gameStreamsCount, len(items))

			for i, item := range items {
				assert.Equal(t, streams[i].UserName, item.Title)
				assert.Equal(
					t,
					fmt.Sprintf("%d viewers - %s", streams[i].ViewerCount, streams[i].Title),
					item.SubTitle,
				)
				assert.Equal(t, fmt.Sprintf("https://www.twitch.tv/%s", streams[i].UserLogin), item.Arg)
			}
		})
	}
}
