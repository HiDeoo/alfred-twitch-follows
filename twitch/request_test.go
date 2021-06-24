package twitch

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetToReturnError(t *testing.T) {
	statusCode := 500

	mockClient := new(MockClient)
	Client = mockClient
	mockClient.On("Do", mock.Anything).Return(mockResponse(statusCode, ""), nil)

	data, err := get("users")

	assert.Nil(t, data)
	assert.EqualValues(t, fmt.Sprintf("Unable to fetch Twitch data (error %d)", statusCode), err.Error())
}

func TestGetToReturnTwitchError(t *testing.T) {
	mockClient := new(MockClient)
	Client = mockClient
	mockClient.On("Do", mock.Anything).Return(mockResponse(
		401,
		`{ "error": "Unauthorized", "status": 401, "message": "OAuth token is missing" }`,
	), nil)

	data, err := get("users")

	assert.Nil(t, data)
	assert.EqualValues(t, "OAuth token is missing", err.Error())
}

func TestGetToReturnData(t *testing.T) {
	json := `{ "data": "the data" }`

	mockClient := new(MockClient)
	Client = mockClient
	mockClient.On("Do", mock.Anything).Return(mockResponse(200, json), nil)

	data, err := get("users")

	assert.EqualValues(t, json, data)
	assert.Nil(t, err)
}
