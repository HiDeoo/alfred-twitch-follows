package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type apiTest struct {
	name                string
	pageCount           int
	lastPageEntityCount int
}

type gqlApiTest struct {
	name  string
	count int
}

func TestGetCurrentUserRepos(t *testing.T) {
	tests := []apiTest{
		{"ReturnNoRepos", 1, 0},
		{"ReturnRepos", 1, 3},
		{"ReturnReposWithPagination", 2, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testAPI(t, test, func(index int) interface{} {
				return GHRepo{
					ID:       123456789 + index,
					FullName: fmt.Sprintf("Full name %d", index),
					HtmlURL:  fmt.Sprintf("https://github.com/user/repo%d", index),
					PushedAt: time.Now().Format(time.RFC3339),
				}
			}, func() (interface{}, error) {
				return GetCurrentUserRepos()
			})
		})
	}
}

func TestGetCurrentUserContributions(t *testing.T) {
	tests := []gqlApiTest{
		{"ReturnNoContributions", 0},
		{"ReturnContributions", 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testGqlAPI(t, test, func(index int) interface{} {
				return GHPullRequestContributionsByRepository{
					Repository: GHRepository{
						ID:            123456789 + index,
						IsFork:        false,
						NameWithOwner: fmt.Sprintf("user/repo %d", index),
						Owner: GHOwner{
							Login: fmt.Sprintf("Login %d", index),
						},
						URL: fmt.Sprintf("https://github.com/user/repo%d", index),
					},
					Contributions: GHContribution{
						Nodes: []GHNode{
							{
								PullRequest: GHPullRequest{
									CreatedAt: time.Now().Format(time.RFC3339),
								},
							},
						},
					},
				}
			}, func() (interface{}, error) {
				return GetCurrentUserContributionRepos()
			})
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
		{"ReturnError", 401, errStr, "unable to fetch GitHub data (error: 401)"},
		{"ReturnGHError", 401, fmt.Sprintf(`{ "message": "%s" , "documentation_url": "http//github.com" }`, errStr), errStr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(request.MockClient)
			client.SetClient(mockClient)

			mockClient.On("Do", mock.Anything).Return(request.MockResponse(test.code, test.response), nil).Once()

			repos, err := GetCurrentUserRepos()

			assert.Nil(t, repos)
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
				string(entityPageJSON),
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

func testGqlAPI(
	t *testing.T,
	test gqlApiTest,
	newEntity func(index int) interface{},
	getEntities func() (interface{}, error),
) {
	t.Run(test.name, func(t *testing.T) {
		mockClient := new(request.MockClient)
		client.SetClient(mockClient)

		entityPage := make([]interface{}, test.count)

		for i := 0; i < test.count; i++ {
			entityJSON := newEntity(i)
			entityPage[i] = entityJSON
		}

		entityPageJSON, err := json.Marshal(entityPage)
		fmt.Println(string(entityPageJSON))
		assert.Nil(t, err)

		mockClient.On("Do", mock.Anything).Return(request.MockResponse(
			200,
			fmt.Sprintf(`{
				"data": {
					"viewer": {
						"contributionsCollection": {
							"pullRequestContributionsByRepository": %s
						}
					}
				}
			}`, entityPageJSON),
		), nil).Once()

		entities, err := getEntities()
		entitySlice, ok := getEntitySlice(entities)

		assert.Equal(t, true, ok)

		assert.Equal(t, test.count, entitySlice.Len())
		assert.Nil(t, err)

		expectedEntities := make([]interface{}, 0)

		for _, entityPage := range entityPage {
			expectedEntities = append(expectedEntities, GHRepo{
				ID:       entityPage.(GHPullRequestContributionsByRepository).Repository.ID,
				FullName: entityPage.(GHPullRequestContributionsByRepository).Repository.NameWithOwner,
				HtmlURL:  entityPage.(GHPullRequestContributionsByRepository).Repository.URL,
				PushedAt: entityPage.(GHPullRequestContributionsByRepository).Contributions.Nodes[0].PullRequest.CreatedAt,
			})
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
