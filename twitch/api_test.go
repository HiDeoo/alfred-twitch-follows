package twitch

import (
	"fmt"
	"strconv"
	"testing"
)

const userJson = `{
	"data":[
		{
			"id": "123456789",
			"login": "user",
			"display_name": "User",
			"type": "",
			"broadcaster_type": "",
			"description": "",
			"profile_image_url": "https://example.com/123456789-profile_image-300x300.png",
			"offline_image_url": "",
			"view_count": 555,
			"email": "user@example.com",
			"created_at": "2020-03-01T02:23:44.009756Z"
		}
	]
}`

func TestGetFollows(t *testing.T) {
	tests := []ApiTest{
		{"ReturnNoFollows", 0, 1, userJson, true},
		{"ReturnFollows", 2, 1, userJson, true},
		{"ReturnFollowsWithPagination", 2, 2, userJson, true},
		{"ReturnError", 2, 1, `{ "data": [] }`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return Follow{
					FromId:     "123456789",
					FromLogin:  "user",
					ToId:       fmt.Sprintf("1234560%d", index),
					ToLogin:    fmt.Sprintf("user%d", index),
					ToName:     fmt.Sprintf("User%d", index),
					FollowedAt: fmt.Sprintf("2020-03-01T02:23:4%d.009756Z", index),
				}
			}, func() (interface{}, error) {
				return GetFollows()
			})
		})
	}
}

func TestGetFollowedStreams(t *testing.T) {
	tests := []ApiTest{
		{"ReturnNoFollowedStreams", 0, 1, userJson, true},
		{"ReturnFollowedStreams", 2, 1, userJson, true},
		{"ReturnFollowedStreamsWithPagination", 2, 2, userJson, true},
		{"ReturnError", 2, 1, `{ "data": [] }`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return Stream{
					Id:           strconv.Itoa(index),
					UserId:       fmt.Sprintf("1234560%d", index),
					UserLogin:    fmt.Sprintf("user%d", index),
					UserName:     fmt.Sprintf("User%d", index),
					GameId:       fmt.Sprintf("game%d", index),
					GameName:     fmt.Sprintf("Game%d", index),
					Type:         "live",
					Title:        fmt.Sprintf("Stream Title %d", index),
					ViewerCount:  index * 100,
					StartedAt:    fmt.Sprintf("2020-02-01T02:23:4%d.009756Z", index),
					Language:     "en",
					ThumbnailUrl: fmt.Sprintf("https://twitch.tv/stream/%d/thumbnail.png", index),
					TagIds:       []string{"1122334455"},
					IsMature:     false,
				}
			}, func() (interface{}, error) {
				return GetFollowedStreams()
			})
		})
	}
}
