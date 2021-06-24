package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type QueryParam struct {
	Key   string
	Value string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

func init() {
	Client = &http.Client{
		Timeout: time.Second * 5,
	}
}

func get(path string, queryParams ...*QueryParam) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiBaseURL+path, nil)

	if err != nil {
		return nil, err
	}

	if len(queryParams) > 0 {
		query := req.URL.Query()

		for _, queryParam := range queryParams {
			query.Add(queryParam.Key, queryParam.Value)
		}

		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TWITCH_OAUTH_TOKEN"))

	res, err := Client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		twitchError := &Error{}

		if err := json.Unmarshal(body, &twitchError); err != nil || twitchError.Error == "" {
			return nil, fmt.Errorf("Unable to fetch Twitch data (error %d)", res.StatusCode)
		}

		return nil, errors.New(twitchError.Message)
	}

	return body, nil
}
