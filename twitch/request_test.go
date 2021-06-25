package twitch

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var tests = []struct {
	name        string
	statusCode  int
	in          string
	success     bool
	out         string
	clientError bool
}{
	{
		"ReturnClientError",
		500,
		"",
		false,
		"Client error",
		true,
	},
	{
		"ReturnNetworkError",
		500,
		"",
		false,
		"Unable to fetch Twitch data (error 500)",
		false,
	},
	{
		"ReturnTwitchError",
		401,
		`{ "error": "Unauthorized", "status": 401, "message": "OAuth token is missing" }`,
		false,
		"OAuth token is missing",
		false,
	},
	{
		"ReturnData",
		200,
		`{ "data": "the data" }`,
		true,
		"",
		false,
	},
}

func TestGet(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(MockClient)
			Client = mockClient
			if test.clientError {
				mockClient.On("Do", mock.Anything).Return(nil, errors.New("Client error"))
			} else {
				mockClient.On("Do", mock.Anything).Return(mockResponse(test.statusCode, test.in), nil)
			}

			data, err := get("users")

			if test.success {
				assert.EqualValues(t, test.in, data)
				assert.Nil(t, err)
			} else {
				assert.Nil(t, data)
				assert.EqualValues(t, test.out, err.Error())
			}
		})
	}
}
