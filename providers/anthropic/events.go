package anthropic

import (
	"encoding/json"

	"github.com/y0ug/llmhaven/chat"
)

// AnthropicEventHandler processes Anthropic-specific events
type AnthropicEventHandler struct {
	message Message
}

func NewAnthropicEventHandler() *AnthropicEventHandler {
	return &AnthropicEventHandler{}
}

func (h *AnthropicEventHandler) ShouldContinue(event MessageStreamEvent) bool {
	return true // event.Type != "message_stop"
}

func (h *AnthropicEventHandler) HandleEvent(
	event MessageStreamEvent,
) (chat.EventStream, error) {
	h.message.Accumulate(event)
	evt := chat.EventStream{Type: event.Type}

	switch event.Type {
	case "error":
		evt.Type = "error"
		evt.Delta = nil
		evt.Message = AnthropicMessageToChatMessage(&h.message)

	case "content_block_delta":
		var delta chat.MessageContent
		if err := json.Unmarshal(event.Delta, &delta); err != nil {
			return evt, nil
		}
		if delta.Type == "text_delta" {
			evt.Type = "text_delta"
			evt.Delta = delta.Text
		}
	case "message_stop":
		evt.Message = AnthropicMessageToChatMessage(&h.message)
	}
	return evt, nil
}
