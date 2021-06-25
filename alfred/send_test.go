package alfred

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
