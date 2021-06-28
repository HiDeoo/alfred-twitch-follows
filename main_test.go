package main

import (
	"fmt"
	"testing"

	"github.com/HiDeoo/alfred-twitch-follows/twitch"
	"github.com/stretchr/testify/assert"
)

func TestMapFollowsToItems(t *testing.T) {
	tests := []struct {
		name        string
		followCount int
	}{
		{"ReturnNoFollows", 0},
		{"ReturnMultipleFollows", 1},
		{"ReturnMultipleFollows", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var follows = []twitch.Follow{}

			for i := 0; i < test.followCount; i++ {
				followJson := twitch.Follow{
					FromId:     "123456789",
					FromLogin:  "user",
					ToId:       fmt.Sprintf("1234560%d", i),
					ToLogin:    fmt.Sprintf("user%d", i),
					ToName:     fmt.Sprintf("User%d", i),
					FollowedAt: fmt.Sprintf("2020-03-01T02:23:4%d.009756Z", i),
				}

				follows = append(follows, followJson)
			}

			items := mapFollowsToItems(follows)

			assert.Equal(t, test.followCount, len(items))

			for i, item := range items {
				assert.Equal(t, follows[i].ToName, item.Title)
				assert.Equal(t, fmt.Sprintf("https://www.twitch.tv/%s", follows[i].ToLogin), item.SubTitle)
				assert.Equal(t, item.SubTitle, item.Arg)
			}
		})
	}
}
