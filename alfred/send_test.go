package alfred

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendResult(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
	}{
		{
			"ReturnItems",
			[]Item{
				{
					"title item 1",
					"subtitle item 1",
					"arg item 1",
				},
				{
					"title item 2",
					"subtitle item 2",
					"arg item 2",
				},
			},
		},
		{
			"ReturnItem",
			[]Item{
				{
					"title single item",
					"subtitle single item",
					"arg single item",
				},
			},
		},
		{
			"ReturnNoItems",
			[]Item{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(func() {
				SendResult(test.items)
			})

			result := Result{}
			err := json.Unmarshal([]byte(output), &result)

			assert.Nil(t, err)

			if len(test.items) > 0 {
				assert.Equal(t, len(test.items), len(result.Items))

				for i, item := range test.items {
					resultItem := result.Items[i].(map[string]interface{})

					assert.Equal(t, item.Title, resultItem["title"])
					assert.Equal(t, item.SubTitle, resultItem["subtitle"])
					assert.Equal(t, item.Arg, resultItem["arg"])
				}
			} else {
				assert.Equal(t, 1, len(result.Items))

				item := result.Items[0].(map[string]interface{})
				placeholder := newEmptyPlaceholderItem()

				assert.Equal(t, placeholder.Title, item["title"])
				assert.Equal(t, placeholder.SubTitle, item["subtitle"])
				assert.Equal(t, placeholder.Arg, item["arg"])
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
