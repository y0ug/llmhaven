package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/config"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
	"github.com/y0ug/llmhaven/internal"
)

// ChatCompletionService implements llmclient.ChatService using OpenAI's types.
type MessageService struct {
	*internal.GenericChatService[MessageNewParams, Message, MessageStreamEvent]
}

func NewMessageService(opts ...options.RequestOption) *MessageService {
	baseService := &internal.GenericChatService[MessageNewParams, Message, MessageStreamEvent]{
		Options:  opts,
		NewError: NewError,
		Endpoint: "v1/messages",
	}

	return &MessageService{
		GenericChatService: baseService,
	}
}

func (svc *MessageService) NewStreaming(
	ctx context.Context,
	params MessageNewParams,
	opts ...options.RequestOption,
) (streaming.Streamer[MessageStreamEvent], error) {
	combinedOpts := append(svc.Options, opts...)
	combinedOpts = append(
		[]options.RequestOption{options.WithJSONSet("stream", true)},
		combinedOpts...)
	path := svc.Endpoint

	var raw *http.Response
	err := config.ExecuteNewRequest(
		ctx,
		http.MethodPost,
		path,
		params,
		&raw,
		svc.NewError,
		combinedOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing new request streaming: %w", err)
	}
	return streaming.NewStream(
		streaming.NewDecoderSSE(raw),
		NewAnthropicStreamHandler(),
	), nil
}

func (svc *MessageService) CountTokens(
	ctx context.Context,
	body MessageNewParams,
	opts ...options.RequestOption,
) (res *MessageTokensCount, err error) {
	opts = append(svc.Options[:], opts...)
	path := "v1/messages/count_tokens"
	err = config.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, svc.NewError, opts...)
	return
}

type MessageTokensCount struct {
	// The total number of tokens across the provided list of messages, system prompt,
	// and tools.
	InputTokens int64 `json:"input_tokens,required"`
}

type MessageParam struct {
	Role    string                 `json:"role"`
	Content []*chat.MessageContent `json:"content"`
}

// Message response, ToParam methode convert to MessageParam
type Message struct {
	ID           string                 `json:"id,omitempty"`
	Content      []*chat.MessageContent `json:"content,omitempty"`
	Role         string                 `json:"role,omitempty"` // Always "assistant"
	StopReason   string                 `json:"stop_reason,omitempty"`
	StopSequence string                 `json:"stop_sequence,omitempty"`
	Type         string                 `json:"type,omitempty"` // Always "message"
	Usage        *Usage                 `json:"usage,omitempty"`
	Model        string                 `json:"model,omitempty"`
}

func (r *Message) ToParam() MessageParam {
	return MessageParam{
		Role:    r.Role,
		Content: r.Content,
	}
}

func (a *Message) Accumulate(event MessageStreamEvent) error {
	if a == nil {
		*a = Message{}
	}

	switch event.Type {
	case "message_start":
		*a = event.Message
	case "content_block_start":
		index := event.Index
		a.Content = append(a.Content, &chat.MessageContent{})
		if int(index) >= len(a.Content) {
			return fmt.Errorf("Index %d is out of range, len: %d\n", index, len(a.Content))
		}

		err := json.Unmarshal(event.ContentBlock, a.Content[index])
		if err != nil {
			return err
		}
	case "content_block_delta":
		index := event.Index
		if int(index) >= len(a.Content) {
			return fmt.Errorf("Index %d is out of range, len: %d\n", index, len(a.Content))
		}
		var delta chat.MessageContent
		err := json.Unmarshal(event.Delta, &delta)
		if err != nil {
			return fmt.Errorf("error unmarshalling delta: %w %s", err, event.Delta)
		}

		switch delta.Type {
		case "text_delta":
			a.Content[index].Text += delta.Text
		case "input_json_delta":
			a.Content[index].InputJson = append(
				a.Content[index].InputJson,
				[]byte(delta.PartialJson)...)
		}
	case "message_delta":

		var delta struct {
			StopReason   string `json:"stop_reason,omitempty"`
			StopSequence string `json:"stop_sequence,omitempty"`
		}
		err := json.Unmarshal(event.Delta, &delta)
		if err != nil {
			return fmt.Errorf("error unmarshalling delta: %w %s", err, event.Delta)
		}
		a.StopReason = delta.StopReason
		a.StopSequence = delta.StopSequence
		a.Usage.OutputTokens = event.Usage.OutputTokens

	//  update StopRead, StopSequence, Usage
	// a.StopReason = event.Delta.StopReason

	case "content_block_stop":
		index := event.Index
		if len(a.Content[index].InputJson) > 0 {
			json.Unmarshal([]byte(a.Content[index].InputJson), &a.Content[index].Input)
		}

	case "message_stop":
		// We should notify the it's complete
	}

	return nil
}

type MessageStreamEvent struct {
	Type string `json:"type"`
	// This field can have the runtime type of [ContentBlockStartEventContentBlock].
	ContentBlock json.RawMessage `json:"content_block"`
	// This field can have the runtime type of [MessageDeltaEventDelta],
	// [ContentBlockDeltaEventDelta].
	Delta   json.RawMessage `json:"delta"`
	Index   int64           `json:"index"`
	Message Message         `json:"message"`
	// Billing and rate-limit usage.
	//
	// Anthropic's API bills and rate-limits by token counts, as tokens represent the
	// underlying cost to our systems.
	//
	// Under the hood, the API transforms requests into a format suitable for the
	// model. The model's output then goes through a parsing stage before becoming an
	// API response. As a result, the token counts in `usage` will not match one-to-one
	// with the exact visible content of an API request or response.
	//
	// For example, `output_tokens` will be non-zero, even for an empty string response
	// from Claude.
	Usage MessageDeltaUsage `json:"usage"`
}

type MessageDeltaUsage struct {
	OutputTokens int `json:"output_tokens"`
}

type Usage struct {
	InputTokens              int `json:"input_tokens,omitempty"`
	OutputTokens             int `json:"output_tokens,omitempty"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

type MessageNewParams struct {
	MaxTokens     int            `json:"max_tokens,omitempty"`
	Messages      []MessageParam `json:"messages"` // MessageParam
	Model         string         `json:"model"`
	StopSequences []string       `json:"stop_sequences,omitempty"`
	Stream        bool           `json:"stream,omitempty"`
	System        string         `json:"system,omitempty"`

	Temperature float64     `json:"temperature,omitempty"` // Number between 0 and 1 that controls randomness of the output.
	Tools       []chat.Tool `json:"tools,omitempty"`       // ToolParam
	ToolChoice  interface{} `json:"tool_choice,omitempty"` // Auto but can be used to force to used a tools
}
