package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type apiTest struct {
	name               string
	entityCountPerPage int
	pageCount          int
	user               string
}

const userJSON = `{
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
}`

var errStr = "Something went wrong!"

func TestGetFollows(t *testing.T) {
	tests := []apiTest{
		{"ReturnNoFollows", 0, 1, userJSON},
		{"ReturnFollows", 2, 1, userJSON},
		{"ReturnFollowsWithPagination", 2, 2, userJSON},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return TwitchFollow{
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
	tests := []apiTest{
		{"ReturnNoFollowedStreams", 0, 1, userJSON},
		{"ReturnFollowedStreams", 2, 1, userJSON},
		{"ReturnFollowedStreamsWithPagination", 2, 2, userJSON},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return TwitchStream{
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

func TestGetCurrentUser(t *testing.T) {
	tests := []struct {
		name      string
		userCount int
	}{
		{"ReturnUser", 1},
		{"ReturnNoUser", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			queryJSON := fmt.Sprintf(`{ "data": [%s] }`, userJSON)

			if test.userCount == 0 {
				queryJSON = `{ "data": [] }`
			}

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(200, queryJSON), nil).Once()

			user, err := getCurrentUser()

			if test.userCount > 0 {
				assert.Nil(t, err)

				expectedUser := TwitchUsers{}

				err = json.Unmarshal([]byte(queryJSON), &expectedUser)
				assert.Nil(t, err)

				assert.Equal(t, &expectedUser.Data[0], user)
			} else {
				assert.Nil(t, user)
				assert.Equal(t, "Unable to get current user", err.Error())
			}
		})
	}
}

func TestQuery(t *testing.T) {
	var errStr = "Something went wrong!"

	tests := []struct {
		name     string
		code     int
		response string
		error    string
	}{
		{"ReturnError", 401, errStr, "Unable to fetch Twitch data (error: 401)"},
		{"ReturnTwitchError", 401, fmt.Sprintf(`{ "error": "Error", "status": 401, "message": "%s" }`, errStr), errStr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(test.code, test.response), nil).Once()

			follows, err := GetFollows()

			assert.Nil(t, follows)
			assert.Equal(t, test.error, err.Error())
		})
	}
}

func testAPI(
	t *testing.T,
	test apiTest,
	newEntity func(index int) interface{},
	getEntities func() (interface{}, error),
) {
	t.Run(test.name, func(t *testing.T) {
		mockClient := new(request.MockClient)
		client.SetClient(mockClient)

		mockClient.On("Do", mock.Anything).Return(
			request.MockResponse(200, fmt.Sprintf(`{ "data": [%s] }`, test.user)),
			nil,
		).Once()

		var cursor string
		entityIndex := 0
		entityPages := make([][]interface{}, test.pageCount)

		for i := 0; i < test.pageCount; i++ {
			entityPages[i] = make([]interface{}, test.entityCountPerPage)

			for j := 0; j < test.entityCountPerPage; j++ {
				entityJSON := newEntity(entityIndex)
				entityPages[i][j] = entityJSON
				entityIndex++
			}

			entityPageJSON, err := json.Marshal(entityPages[i])
			assert.Nil(t, err)

			if test.pageCount > 1 && i != test.pageCount-1 {
				cursor = "abcdefgh"
			} else {
				cursor = ""
			}

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(
				200,
				fmt.Sprintf(`{ "data": %s, "pagination": { "cursor": "%s" } }`, entityPageJSON, cursor),
			), nil).Once()
		}

		entities, err := getEntities()
		entitySlice, ok := getEntitySlice(entities)

		assert.Equal(t, true, ok)

		assert.Equal(t, test.entityCountPerPage*test.pageCount, entitySlice.Len())
		assert.Nil(t, err)

		expectedEntities := make([]interface{}, 0)

		for _, entityPage := range entityPages {
			expectedEntities = append(expectedEntities, entityPage...)
		}

		assert.ElementsMatch(t, expectedEntities, entities)
	})
}

func getEntitySlice(entities interface{}) (reflect.Value, bool) {
	ok := false
	val := reflect.ValueOf(entities)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return val, ok
}
