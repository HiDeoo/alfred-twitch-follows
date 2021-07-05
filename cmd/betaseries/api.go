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
	queryLimit = 199
)

func init() {
	client = request.NewClient("https://api.betaseries.com/")
	client.SetHeaders(http.Header{
		"X-BetaSeries-Version": []string{"3.0"},
		"X-BetaSeries-Key":     []string{os.Getenv("BETASERIES_CLIENT_ID")},
		"Authorization":        []string{"Bearer " + os.Getenv("BETASERIES_OAUTH_TOKEN")},
	})
}

func getCurrentShowsWithUnwatchedEpisodes() ([]BSShow, error) {
	var offset = 0
	var allShows []BSShow

	for {
		shows, err := getCurrentShowsWithUnwatchedEpisodesAndPagination(offset)

		if err != nil {
			return nil, err
		}

		allShows = append(allShows, shows.Shows...)

		if len(shows.Shows) < queryLimit || len(shows.Shows) == 0 {
			break
		}

		offset = offset + queryLimit
	}

	return allShows, nil
}

func getCurrentShowsWithUnwatchedEpisodesAndPagination(offset int) (*BSShows, error) {
	queryParams := url.Values{}
	queryParams.Set("status", "current")
	queryParams.Set("offset", strconv.Itoa(offset))
	queryParams.Set("limit", strconv.Itoa(queryLimit))

	res, err := query(client.Get("shows/member", queryParams))

	if err != nil {
		return nil, err
	}

	shows := BSShows{}

	if err = json.Unmarshal(res.Data, &shows); err != nil {
		return nil, err
	}

	return &shows, nil
}

func query(res *request.Response, err error) (*request.Response, error) {
	if res.StatusCode != 200 {
		bsError := &BSErrors{}

		json.Unmarshal(res.Data, &bsError)

		if err := json.Unmarshal(res.Data, &bsError); err != nil {
			return nil, fmt.Errorf("Unable to fetch BetaSeries data (error: %d)", res.StatusCode)
		}

		return nil, errors.New(bsError.Errors[0].Text)
	}

	return res, err
}
