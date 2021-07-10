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
	name                string
	pageCount           int
	lastPageEntityCount int
}

func TestGetCurrentShowsWithUnwatchedEpisodes(t *testing.T) {
	tests := []apiTest{
		{"ReturnNoShows", 1, 0},
		{"ReturnShows", 1, 3},
		{"ReturnShowsWithPagination", 2, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return BSShow{
					ID:          123456789 + index,
					Title:       fmt.Sprintf("Title %d", index),
					Description: fmt.Sprintf("Description %d", index),
					Episodes:    strconv.Itoa(100 + index),
				}
			}, func() (interface{}, error) {
				return getCurrentShowsWithUnwatchedEpisodes()
			})
		})
	}
}

func TestMarkShowAsWatched(t *testing.T) {
	tests := []struct {
		name    string
		success bool
	}{
		{"ReturnNoError", true},
		{"ReturnError", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			if test.success {
				mockClient.On("Do", mock.Anything).Return(request.MockResponse(200, `{ "episode": "data" }`), nil).Once()
			} else {
				mockClient.On("Do", mock.Anything).Return(
					request.MockResponse(400, `{ "errors": [{ "code": 4001 , "text": "L'épisode avec l'ID X n'existe pas." }] }`),
					nil,
				).Once()
			}

			err := markEpisodeAsWatched(987654321)

			if test.success {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestMarkEpisodeAsWatched(t *testing.T) {
	tests := []struct {
		name                 string
		success              bool
		lastAiredEpisodeDate string
	}{
		{"ReturnNoError", true, "2020-10-22"},
		{"ReturnErrorWithNoLastAiredEpisode", false, "2025-11-22"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			episodes := BSEpisodes{[]BSEpisode{{ID: 987654321, Date: "2020-10-22"}}, BSErrors{}}
			episodesJSON, err := json.Marshal(episodes)
			assert.Nil(t, err)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(200, string(episodesJSON)), nil).Once()

			if test.success {
				mockClient.On("Do", mock.Anything).Return(request.MockResponse(200, `{ "episode": "data" }`), nil).Once()
			} else {
				mockClient.On("Do", mock.Anything).Return(
					request.MockResponse(400, `{ "errors": [{ "code": 4001 , "text": "L'épisode avec l'ID X n'existe pas." }] }`),
					nil,
				).Once()
			}

			err = markShowAsWatched("123456789")

			if test.success {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetLastAiredEpisode(t *testing.T) {
	tests := []struct {
		name             string
		lastEpisodeIndex int
		episodeDates     []string
	}{
		{"ReturnLastEpisodeAtLastIndex", 1, []string{"2020-10-22", "2020-11-22"}},
		{"ReturnLastEpisodeNotAtLastIndex", 2, []string{"2020-10-22", "2020-11-22", "2020-12-22", "2025-11-22"}},
		{"ReturnNoLastEpisode", -1, []string{"2025-10-22", "2025-11-22"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			episodes := make([]BSEpisode, len(test.episodeDates))

			for i, date := range test.episodeDates {
				episodes[i] = BSEpisode{
					Date: date,
				}
			}

			episodesJSON, err := json.Marshal(episodes)
			assert.Nil(t, err)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(
				200,
				fmt.Sprintf(`{ "episodes": %s }`, episodesJSON),
			), nil).Once()

			lastEpisode, err := getLastAiredEpisode("123")

			if test.lastEpisodeIndex != -1 {
				assert.Nil(t, err)
				assert.Equal(t, test.episodeDates[test.lastEpisodeIndex], lastEpisode.Date)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, lastEpisode)
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
		{"ReturnError", 401, errStr, "Unable to fetch BetaSeries data (error: 401)"},
		{"ReturnBSError", 401, fmt.Sprintf(`{ "errors": [{ "code": 401 , "text": "%s" }] }`, errStr), errStr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(test.code, test.response), nil).Once()

			follows, err := getCurrentShowsWithUnwatchedEpisodes()

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

		entityIndex := 0
		entityPages := make([][]interface{}, test.pageCount)

		for i := 0; i < test.pageCount; i++ {
			entityCount := queryLimit

			if i == test.pageCount-1 {
				entityCount = test.lastPageEntityCount
			}

			entityPages[i] = make([]interface{}, entityCount)

			for j := 0; j < entityCount; j++ {
				entityJSON := newEntity(entityIndex)
				entityPages[i][j] = entityJSON
				entityIndex++
			}

			entityPageJSON, err := json.Marshal(entityPages[i])
			assert.Nil(t, err)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(
				200,
				fmt.Sprintf(`{ "shows": %s }`, entityPageJSON),
			), nil).Once()
		}

		entities, err := getEntities()
		entitySlice, ok := getEntitySlice(entities)

		assert.Equal(t, true, ok)

		assert.Equal(t, queryLimit*(test.pageCount-1)+test.lastPageEntityCount, entitySlice.Len())
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
