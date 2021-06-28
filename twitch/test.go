package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ApiTest struct {
	name               string
	entityCountPerPage int
	pageCount          int
	user               string
	success            bool
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	arg0 := args.Get(0)

	if arg0 != nil {
		return args.Get(0).(*http.Response), args.Error(1)
	}

	return nil, args.Error(1)
}

func mockResponse(statusCode int, json string) *http.Response {
	readCloser := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	return &http.Response{StatusCode: statusCode, Body: readCloser}
}

func testAPI(
	t *testing.T,
	test ApiTest,
	newEntity func(index int) interface{},
	getEntities func() (interface{}, error),
) {
	t.Run(test.name, func(t *testing.T) {
		mockClient := new(MockClient)
		Client = mockClient
		mockClient.On("Do", mock.Anything).Return(mockResponse(200, test.user), nil).Once()

		var cursor string
		entityIndex := 0
		entityPages := make([][]interface{}, test.pageCount)

		for i := 0; i < test.pageCount; i++ {
			entityPages[i] = make([]interface{}, test.entityCountPerPage)

			for j := 0; j < test.entityCountPerPage; j++ {
				entityJSON := newEntity(entityIndex)
				entityPages[i][j] = entityJSON
				entityIndex++
			}

			entityPageJSON, err := json.Marshal(entityPages[i])
			assert.Nil(t, err)

			if test.pageCount > 1 && i != test.pageCount-1 {
				cursor = "abcdefgh"
			} else {
				cursor = ""
			}

			json := fmt.Sprintf(`{ "data": %s, "pagination": { "cursor": "%s" } }`, entityPageJSON, cursor)

			mockClient.On("Do", mock.Anything).Return(mockResponse(200, json), nil).Once()
		}

		entities, err := getEntities()
		entitySlice, ok := getEntitySlice(entities)

		assert.Equal(t, true, ok)

		if test.success {
			assert.Equal(t, test.entityCountPerPage*test.pageCount, entitySlice.Len())
			assert.Nil(t, err)

			expectedEntities := make([]interface{}, 0)

			for _, entityPage := range entityPages {
				expectedEntities = append(expectedEntities, entityPage...)
			}

			assert.ElementsMatch(t, expectedEntities, entities)
		} else {
			assert.Nil(t, entities)
			assert.NotNil(t, err)
		}
	})
}

func getEntitySlice(entities interface{}) (reflect.Value, bool) {
	ok := false
	val := reflect.ValueOf(entities)

	if val.Kind() == reflect.Slice {
		ok = true
	}

	return val, ok
}
