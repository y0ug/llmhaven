package openai

import "github.com/y0ug/llmhaven/chat"

// OpenAIEventHandler processes OpenAI-specific events
type OpenAIEventHandler struct {
	completion ChatCompletion
}

func NewOpenAIEventHandler() *OpenAIEventHandler {
	return &OpenAIEventHandler{}
}

func (h *OpenAIEventHandler) ShouldContinue(chunk ChatCompletionChunk) bool {
	return true
	// return !(chunk.Usage.CompletionTokens != 0 || len(chunk.Choices) == 0)
}

func (h *OpenAIEventHandler) HandleEvent(
	chunk ChatCompletionChunk,
) (chat.EventStream, error) {
	h.completion.Accumulate(chunk)
	evt := chat.EventStream{Message: ToChatResponse(&h.completion)}

	if chunk.Usage.CompletionTokens != 0 || len(chunk.Choices) == 0 {
		evt.Type = "message_stop"
		return evt, nil
	}

	evt.Type = "text_delta"
	evt.Delta = chunk.Choices[0].Delta.Content
	return evt, nil
}
