package request

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (mockClient *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := mockClient.Called(req)
	arg0 := args.Get(0)

	if arg0 != nil {
		if len(args) >= 3 {
			validator := args.Get(2).(func(req *http.Request))
			validator(req)
		}

		return args.Get(0).(*http.Response), args.Error(1)
	}

	return nil, args.Error(1)
}

func MockResponse(statusCode int, json string) *http.Response {
	readCloser := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	return &http.Response{StatusCode: statusCode, Body: readCloser}
}
