package anthropic

import (
	"encoding/json"
	"fmt"

	"github.com/y0ug/llmhaven/http/streaming"
)

// AnthropicStreamHandler implements BaseStreamHandler for Anthropic's streaming responses
type AnthropicStreamHandler struct{}

func NewAnthropicStreamHandler() *AnthropicStreamHandler {
	return &AnthropicStreamHandler{}
}

func (h *AnthropicStreamHandler) HandleEvent(event streaming.Event) (MessageStreamEvent, error) {
	var result MessageStreamEvent
	var err error
	switch event.Type {
	case "completion":
		if err := json.Unmarshal(event.Data, &result); err != nil {
			return result, err
		}
	case "message_start",
		"message_delta",
		"message_stop",
		"content_block_start",
		"content_block_delta",
		"content_block_stop":
		if err := json.Unmarshal(event.Data, &result); err != nil {
			return result, err
		}
	case "error":
		var errorResp struct {
			Error struct {
				Type    string      `json:"type"`
				Message string      `json:"message"`
				Details interface{} `json:"details"`
			} `json:"error"`
		}
		if err := json.Unmarshal(event.Data, &errorResp); err != nil {
			return result, fmt.Errorf("failed to parse error response: %w", err)
		}
		if err := json.Unmarshal(event.Data, &result); err != nil {
			return result, fmt.Errorf("failed to parse error response: %w", err)
		}
		// fmt.Println("result", string(result.data))
		err = fmt.Errorf("%s", errorResp.Error.Message)
	}

	return result, err
}

func (h *AnthropicStreamHandler) ShouldContinue(event streaming.Event) bool {
	if event.Type == "ping" {
		return true
	}
	return event.Type != "error"
}
