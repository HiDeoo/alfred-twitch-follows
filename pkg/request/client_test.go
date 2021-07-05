package request

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		response    string
		queryParams url.Values
		headers     http.Header
	}{
		{"ReturnClientError", 0, "Client error", url.Values{}, http.Header{}},
		{"ReturnError", 401, `{ "error": "Unauthorized" }`, url.Values{}, http.Header{}},
		{"ReturnData", 200, `{ "data": "the data" }`, url.Values{}, http.Header{}},
		{"ReturnDataWithHeaders", 200, "", url.Values{}, http.Header{"headerKey": {"headerValue"}}},
		{"ReturnDataWithQueryParams", 200, "", url.Values{"queryKey": []string{"queryValue"}}, http.Header{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient("https://example.com")
			mockClient := new(MockClient)
			client.SetClient(mockClient)

			if len(test.headers) > 0 {
				client.SetHeaders(test.headers)
			}

			if test.statusCode == 0 {
				mockClient.On("Do", mock.Anything).Return(nil, errors.New(test.response))
			} else {
				mockClient.On("Do", mock.Anything).Return(
					MockResponse(test.statusCode, test.response),
					nil,
					func(req *http.Request) {

						assert.Equal(t, "example.com", req.URL.Host)
						assert.Equal(t, "/fake", req.URL.Path)

						if len(test.queryParams) > 0 {
							queryParams := req.URL.Query()

							assert.EqualValues(t, test.queryParams, queryParams)
						}

						if len(test.headers) > 0 {
							assert.EqualValues(t, test.headers, req.Header)
						}
					},
				)
			}

			res, err := client.Get("fake", test.queryParams)

			if test.statusCode == 0 {
				assert.Nil(t, res)
				assert.EqualValues(t, test.response, err.Error())
			} else {
				assert.Equal(t, test.statusCode, res.StatusCode)
				assert.Equal(t, []byte(test.response), res.Data)
				assert.Nil(t, err)
			}
		})
	}
}
