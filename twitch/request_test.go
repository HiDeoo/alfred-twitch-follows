package twitch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var tests = []struct {
	name       string
	statusCode int
	in         string
	success    bool
	out        string
}{
	{
		"ReturnError",
		500,
		"",
		false,
		"Unable to fetch Twitch data (error 500)"},
	{
		"ReturnTwitchError",
		401,
		`{ "error": "Unauthorized", "status": 401, "message": "OAuth token is missing" }`,
		false,
		"OAuth token is missing",
	},
	{
		"ReturnData",
		200,
		`{ "data": "the data" }`,
		true,
		"",
	},
}

func TestGet(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(MockClient)
			Client = mockClient
			mockClient.On("Do", mock.Anything).Return(mockResponse(test.statusCode, test.in), nil)

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
