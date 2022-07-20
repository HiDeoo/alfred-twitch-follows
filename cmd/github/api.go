package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/request"
)

var (
	client     request.Client
	queryLimit = 100
)

func init() {
	client = request.NewClient("https://api.github.com/")
	client.SetHeaders(http.Header{
		"Accept":        []string{"application/vnd.github+json"},
		"Authorization": []string{"Token " + os.Getenv("GITHUB_OAUTH_TOKEN")},
	})
}

func GetCurrentUserRepos() ([]GHRepo, error) {
	var page = 1
	var allRepos []GHRepo

	for {
		repos, err := GetCurrentUserReposWithPagination(page)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < queryLimit || len(repos) == 0 {
			break
		}

		page = page + 1
	}

	return allRepos, nil
}

func GetCurrentUserReposWithPagination(page int) ([]GHRepo, error) {
	queryParams := url.Values{}
	queryParams.Set("sort", "pushed")
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("per_page", strconv.Itoa(queryLimit))

	res, err := query(client.Get("user/repos", queryParams))

	if err != nil {
		return nil, err
	}

	repos := []GHRepo{}

	if err = json.Unmarshal(res.Data, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func query(res *request.Response, err error) (*request.Response, error) {
	if res.StatusCode != 200 {
		ghError := &GHError{}

		if err := json.Unmarshal(res.Data, &ghError); err != nil || ghError.Message == "" {
			return nil, fmt.Errorf("unable to fetch GitHub data (error: %d)", res.StatusCode)
		}

		return nil, errors.New(ghError.Message)
	}

	return res, err
}
