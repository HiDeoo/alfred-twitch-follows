package twitch

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func mockResponse(statusCode int, json string) *http.Response {
	readCloser := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	return &http.Response{StatusCode: statusCode, Body: readCloser}
}
