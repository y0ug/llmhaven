package chat

import (
	"context"

	"github.com/y0ug/llmhaven/http/streaming"
)

type ChatParams struct {
	Model       string
	MaxTokens   int
	Temperature float64
	Messages    []*ChatMessage
	Stream      bool
	Tools       []Tool
	ToolChoice  string // auto, any, tool
	N           *int   // number of choice
}

type ChatResponse struct {
	ID     string       `json:"id,omitempty"`
	Choice []ChatChoice `json:"choice,omitempty"`
	Usage  *ChatUsage   `json:"usage,omitempty"`
	Model  string       `json:"model,omitempty"`
}

func (cm *ChatResponse) ToMessageParams() *ChatMessage {
	if len(cm.Choice) == 0 {
		return nil
	}
	return &ChatMessage{
		Content: cm.Choice[0].Content,
		Role:    cm.Choice[0].Role,
	}
}

func (cm *ChatResponse) HasContent() bool {
	if len(cm.Choice) == 0 {
		return false
	}
	if len(cm.Choice[0].Content) == 0 {
		return false
	}
	return true
}

// StopReason is the reason the model stopped generating messages. It can be one of:
// - `"end_turn"`: the model reached a natural stopping point
// - `"max_tokens"`: we exceeded the requested `max_tokens` or the model's maximum
// - `"stop_sequence"`: one of your provided custom `stop_sequences` was generated
// - `"tool_use"`: the model invoked one or more tools
type ChatChoice struct {
	Role       string            `json:"role,omitempty"` // Always "assistant"
	Content    []*MessageContent `json:"content,omitempty"`
	StopReason string            `json:"stop_reason,omitempty"`
}

type ChatUsage struct {
	OutputTokens             int `json:"output_tokens"`
	OutputAudioTokens        int `json:"output_audio_tokens"`
	OutputReasoningTokens    int `json:"output_reasoning_tokens"`
	InputTokens              int `json:"input_tokens"`
	InputAudioTokens         int `json:"input_audio_tokens"`
	InputCachedTokens        int `json:"input_cached_tokens"`
	InputCacheCreationTokens int `json:"input_cache_creation_tokens"`
}

type ChatMessage struct {
	Role    string            `json:"role"`
	Content []*MessageContent `json:"content"`
}

func (cm *ChatMessage) SetCache() {
	for _, c := range cm.Content {
		c.SetCache()
	}
}

type Tool struct {
	Description *string     `json:"description,omitempty"`
	InputSchema interface{} `json:"input_schema,omitempty"`
	Name        string      `json:"name"`
}

func StreamChatMessageToChannel(
	ctx context.Context,
	stream streaming.Streamer[EventStream],
	ch chan<- EventStream,
) error {
	defer close(ch)

	for stream.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ch <- stream.Current()
		}
	}

	return stream.Err()
}
