package alfred

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sendResultTests = []struct {
	name  string
	items []Item
}{
	{"ReturnItems", []Item{{"title item 1", "subtitle item 1"}, {"title item 2", "subtitle item 2"}}},
	{"ReturnItem", []Item{{"title single item", "subtitle single item"}}},
	{"ReturnNoItems", []Item{}},
}

func TestSendResult(t *testing.T) {
	for _, test := range sendResultTests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(func() {
				SendResult(test.items)
			})

			result := Result{}
			err := json.Unmarshal([]byte(output), &result)

			assert.Nil(t, err)

			if len(test.items) > 0 {
				for i, item := range test.items {
					resultItem := result.Items[i].(map[string]interface{})

					assert.Equal(t, resultItem["title"], item.Title)
					assert.Equal(t, resultItem["subtitle"], item.SubTitle)
				}
			} else {
				// TODO(HiDeoo)
			}
		})
	}
}

func TestSendError(t *testing.T) {
	errStr := "test error"

	output := captureOutput(func() {
		SendError(errors.New(errStr))
	})

	result := Result{}
	err := json.Unmarshal([]byte(output), &result)

	assert.Equal(t, 1, len(result.Items))
	assert.Nil(t, err)

	item := result.Items[0].(map[string]interface{})

	assert.Equal(t, item["title"], "Something went wrong!")
	assert.Equal(t, item["subtitle"], errStr)
	assert.Equal(t, item["valid"], false)
}
