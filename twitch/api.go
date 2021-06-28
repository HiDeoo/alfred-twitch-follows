package twitch

import (
	"encoding/json"
	"errors"
)

const apiBaseURL = "https://api.twitch.tv/helix/"

func GetFollows() ([]Follow, error) {
	currentUser, err := getCurrentUser()

	if err != nil {
		return nil, err
	}

	var cursor = ""
	var allFollows []Follow

	for {
		follows, err := getFollowsWithPagination(currentUser.Id, cursor)

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

func GetFollowedStreams() ([]Stream, error) {
	currentUser, err := getCurrentUser()

	if err != nil {
		return nil, err
	}

	var cursor = ""
	var allStreams []Stream

	for {
		streams, err := getFollowedStreamsWithPagination(currentUser.Id, cursor)

		if err != nil {
			return nil, err
		}

		allStreams = append(allStreams, streams.Data...)

		if len(streams.Pagination.Cursor) == 0 {
			break
		}

		cursor = streams.Pagination.Cursor
	}

	return allStreams, nil
}

func getFollowsWithPagination(userID string, cursor string) (*Follows, error) {
	data, err := get(
		"users/follows",
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

func getFollowedStreamsWithPagination(userID string, cursor string) (*Streams, error) {
	data, err := get(
		"streams/followed",
		&QueryParam{"user_id", userID},
		&QueryParam{"first", "100"},
		&QueryParam{"after", cursor},
	)

	if err != nil {
		return nil, err
	}

	streams := Streams{}

	if err = json.Unmarshal(data, &streams); err != nil {
		return nil, err
	}

	return &streams, nil
}

func getCurrentUser() (*User, error) {
	data, err := get("users")

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
