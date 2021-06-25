package twitch

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

var getFollowsTests = []struct {
	name        string
	followCount int
	pageCount   int
	user        string
	success     bool
}{
	{"ReturnNoFollows", 0, 1, userJson, true},
	{"ReturnFollows", 2, 1, userJson, true},
	{"ReturnFollowsWithPagination", 2, 2, userJson, true},
	{"ReturnError", 2, 1, `{ "data": [] }`, false},
}

func TestGetFollows(t *testing.T) {
	for _, test := range getFollowsTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(MockClient)
			Client = mockClient
			mockClient.On("Do", mock.Anything).Return(mockResponse(200, test.user), nil).Once()

			var cursor string
			followIndex := 0
			followsPerPage := make([][]Follow, test.pageCount)

			for i := 0; i < test.pageCount; i++ {
				followsPerPage[i] = []Follow{}

				for j := 0; j < test.followCount; j++ {
					followJson := Follow{
						FromId:     "123456789",
						FromLogin:  "user",
						ToId:       fmt.Sprintf("1234560%d", followIndex),
						ToLogin:    fmt.Sprintf("user%d", followIndex),
						ToName:     fmt.Sprintf("User%d", followIndex),
						FollowedAt: fmt.Sprintf("2020-03-01T02:23:4%d.009756Z", followIndex),
					}

					followsPerPage[i] = append(followsPerPage[i], followJson)
					followIndex++
				}

				pageFollowsJson, err := json.Marshal(followsPerPage[i])
				assert.Nil(t, err)

				if test.pageCount > 1 && i != test.pageCount-1 {
					cursor = "abcdefgh"
				} else {
					cursor = ""
				}

				json := fmt.Sprintf(
					`{ "total": %d, "data": %s, "pagination": { "cursor": "%s" } }`,
					len(followsPerPage[i]),
					pageFollowsJson,
					cursor,
				)

				mockClient.On("Do", mock.Anything).Return(mockResponse(200, json), nil).Once()
			}

			follows, err := GetFollows()

			if test.success {
				assert.Equal(t, test.followCount*test.pageCount, len(follows))
				assert.Nil(t, err)

				expectedFollows := []Follow{}

				for _, pageFollows := range followsPerPage {
					expectedFollows = append(expectedFollows, pageFollows...)
				}

				assert.ElementsMatch(t, expectedFollows, follows)
			} else {
				assert.Nil(t, follows)
				assert.NotNil(t, err)
			}
		})
	}
}
