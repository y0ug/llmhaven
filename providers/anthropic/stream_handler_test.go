package anthropic

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/y0ug/llmhaven/http/streaming"
)

func TestAnthropicStreamHandler_HandleEvent_Error(t *testing.T) {
	handler := NewAnthropicStreamHandler()

	// Create an error event similar to the one in the trace
	errorData := map[string]interface{}{
		"type": "error",
		"error": map[string]interface{}{
			"type":    "overloaded_error",
			"message": "Overloaded",
			"details": nil,
		},
	}

	eventData, err := json.Marshal(errorData)
	assert.NoError(t, err)

	event := streaming.Event{
		Type: "error",
		Data: eventData,
	}

	result, err := handler.HandleEvent(event)
	fmt.Println(result)
	fmt.Println(err)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Overloaded")
}
