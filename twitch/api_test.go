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
}{
	{"ReturnNoFollows", 0, 1},
	{"ReturnFollows", 2, 1},
	{"ReturnFollowsWithPagination", 2, 2},
}

func TestGetFollows(t *testing.T) {
	for _, test := range getFollowsTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(MockClient)
			Client = mockClient
			mockClient.On("Do", mock.Anything).Return(mockResponse(200, userJson), nil).Once()

			followsInput := []Follow{}

			for i := 0; i < test.followCount; i++ {
				followJson := Follow{
					FromId:     "123456789",
					FromLogin:  "user",
					ToId:       fmt.Sprintf("1234560%d", i),
					ToLogin:    fmt.Sprintf("user%d", i),
					ToName:     fmt.Sprintf("User%d", i),
					FollowedAt: fmt.Sprintf("2020-03-01T02:23:4%d.009756Z", i),
				}

				followsInput = append(followsInput, followJson)
			}

			followsJson, err := json.Marshal(followsInput)
			assert.Nil(t, err)

			var pagination string

			for i := 0; i < test.pageCount; i++ {
				if test.pageCount > 1 && i == test.pageCount-1 {
					pagination = "abcdefgh"
				}

				json := fmt.Sprintf(
					`{ "total": %d, "data": %s, "pagination": { "cursor": "%s" } }`,
					len(followsInput),
					followsJson,
					pagination,
				)

				mockClient.On("Do", mock.Anything).Return(mockResponse(200, json), nil).Once()
			}

			follows, err := GetFollows()

			assert.Equal(t, test.followCount, len(follows))
			assert.ElementsMatch(t, followsInput, follows)
			assert.Nil(t, err)
		})
	}
}
