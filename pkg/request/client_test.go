package request

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{"Get", http.MethodGet},
		{"Post", http.MethodPost},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient("https://example.com/")
			mockClient := new(MockClient)
			client.SetClient(mockClient)

			queryParams := url.Values{"queryKey": []string{"queryValue"}}

			mockClient.On("Do", mock.Anything).Return(
				MockResponse(200, ""),
				nil,
				func(req *http.Request) {
					assert.Equal(t, test.method, req.Method)

					assert.Equal(t, "/fake", req.URL.Path)

					queryParams := req.URL.Query()
					assert.EqualValues(t, queryParams, queryParams)
				},
			)

			var res *Response
			var err error

			switch test.method {
			case http.MethodGet:
				res, err = client.Get("fake", queryParams)
			case http.MethodPost:
				res, err = client.Post("fake", queryParams)
			}

			assert.NotNil(t, res)
			assert.Nil(t, err)
		})
	}
}

func TestRequest(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		response    string
		queryParams url.Values
		headers     http.Header
		method      string
	}{
		{"GetClientError", 0, "Client error", url.Values{}, http.Header{}, http.MethodGet},
		{"GetError", 401, `{ "error": "Unauthorized" }`, url.Values{}, http.Header{}, http.MethodGet},
		{"GetData", 200, `{ "data": "the data" }`, url.Values{}, http.Header{}, http.MethodGet},
		{"GetDataWithHeaders", 200, "", url.Values{}, http.Header{"headerKey": {"headerValue"}}, http.MethodGet},
		{"GetDataWithQueryParams", 200, "", url.Values{"queryKey": []string{"queryValue"}}, http.Header{}, http.MethodGet},
		{"PostData", 200, `{ "data": "the data" }`, url.Values{}, http.Header{}, http.MethodGet},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := NewClient("https://example.com/")
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
						assert.Equal(t, test.method, req.Method)

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

			res, err := client.request(test.method, "fake", test.queryParams)

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
