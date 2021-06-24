package twitch

import (
	"encoding/json"
	"errors"
	"net/http"
)

const apiBaseURL = "https://api.twitch.tv/helix/"

func GetFollows(httpClient *http.Client) ([]Follow, error) {
	currentUser, err := getCurrentUser(httpClient)

	if err != nil {
		return nil, err
	}

	var cursor = ""
	var allFollows []Follow

	for {
		follows, err := getFollowsWithPagination(httpClient, currentUser.Id, cursor)

		if err != nil {
			return nil, err
		}

		allFollows = append(allFollows, follows.Data...)

		if len(follows.Pagination.Cursor) == 0 {
			break
		}

		cursor = follows.Pagination.Cursor
	}

	return allFollows, nil
}

func getFollowsWithPagination(httpClient *http.Client, userID string, cursor string) (*Follows, error) {
	data, err := get(
		httpClient, "users/follows",
		&QueryParam{"from_id", userID},
		&QueryParam{"first", "100"},
		&QueryParam{"after", cursor},
	)

	if err != nil {
		return nil, err
	}

	follows := Follows{}

	if err = json.Unmarshal(data, &follows); err != nil {
		return nil, err
	}

	return &follows, nil
}

func getCurrentUser(httpClient *http.Client) (*User, error) {
	data, err := get(httpClient, "users")

	if err != nil {
		return nil, err
	}

	users := Users{}

	if err = json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	if len(users.Data) != 1 {
		return nil, errors.New("Unable to get current user")
	}

	return &users.Data[0], nil
}
