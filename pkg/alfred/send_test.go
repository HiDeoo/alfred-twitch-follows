package alfred

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
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
					BaseItem: BaseItem{"title item 1", "subtitle item 1"},
					Arg:      "arg item 1",
				},
				{
					BaseItem: BaseItem{"title item 2", "subtitle item 2"},
					Arg:      "arg item 2",
				},
			},
		},
		{
			"ReturnItem",
			[]Item{
				{
					BaseItem: BaseItem{"title single item", "subtitle single item"},
					Arg:      "arg single item",
				},
			},
		},
		{
			"ReturnItemWithModifiers",
			[]Item{
				{
					BaseItem: BaseItem{"title item", "subtitle item"},
					Arg:      "arg item",
					Mods:     &Modifiers{Cmd: Modifier{Valid: true}},
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

			fmt.Println("output", output)

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

					if item.Mods == nil {
						assert.Nil(t, resultItem["mods"])
					} else {
						assert.NotNil(t, resultItem["mods"])
					}
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

// https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func captureOutput(fn func()) string {
	reader, writer, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr

	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)

	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()

	wg.Wait()

	fn()

	writer.Close()

	return <-out
}
