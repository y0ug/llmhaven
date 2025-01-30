package chat

import (
	"encoding/json"
	"fmt"
	"log"
)

// MessageContentType enumerates possible content types we handle
type MessageContentType string

const (
	// Common Anthropic/OpenAI
	ContentTypeText      MessageContentType = "text"
	ContentTypeTextDelta MessageContentType = "text_delta"

	// Anthropic
	ContentTypeInputJsonDelta MessageContentType = "input_json_delta"
	ContentTypeToolUse        MessageContentType = "tool_use"
	ContentTypeToolResult     MessageContentType = "tool_result"
	ContentTypeDocument       MessageContentType = "document"
	ContentTypeImage          MessageContentType = "image"

	// OpenAI
	ContentTypeInputAudio MessageContentType = "input_audio"
)

// MessageContent holds text, tool calls, or other specialized content
type MessageContent struct {
	Type MessageContentType `json:"type"`

	// Relevant for text content
	Text        string `json:"text,omitempty"`
	PartialJson string `json:"partial_json,omitempty"`

	// Relevant for tool usage calls (like "function calls")
	ID        string          `json:"id,omitempty"`    // Unique identifier for this tool call
	Name      string          `json:"name,omitempty"`  // Name of the tool to call
	Input     json.RawMessage `json:"input,omitempty"` // Arguments to pass to the tool
	InputJson []byte          `json:"-"`               // Arguments to pass to the tool in json format

	// Relevant for tool results
	ToolUseID    string        `json:"tool_use_id,omitempty"`   // ID of the tool call this result is for
	Content      string        `json:"content,omitempty"`       // Result returned from the tool
	Source       *AIContentSrc `json:"source,omitempty"`        // Source of the content if type document/image
	CacheControl *CacheControl `json:"cache_control,omitempty"` // Used to set cache
}

func (m *MessageContent) SetCache() {
	m.CacheControl = NewCacheControlEphemeral()
}

func (m *MessageContent) IsCacheable() bool {
	return m.CacheControl != nil && m.CacheControl.Type == "ephemeral"
}

type CacheControl struct {
	Type string `json:"type"`
}

func NewCacheControlEphemeral() *CacheControl {
	return &CacheControl{
		Type: "ephemeral",
	}
}

type AIContentSrc struct {
	Type      string `json:"type"`       // base64
	MediaType string `json:"media_type"` // "application/pdf" "image/jpeg" etc..
	Data      []byte `json:"data"`
}

func NewSourceContent(sourceType string, mediaType string, data []byte) *MessageContent {
	var contentType MessageContentType
	switch sourceType {
	case "document":
		contentType = ContentTypeDocument
	case "image":
		contentType = ContentTypeImage
	default:
		// TODO: remove this log
		log.Printf("Unknown source type: %s", sourceType)
		contentType = MessageContentType("contentType")
	}
	return &MessageContent{
		Type: contentType,
		Source: &AIContentSrc{
			Type:      "base64",
			MediaType: mediaType,
			Data:      data,
		},
	}
}

// GetType returns the content type
func (c MessageContent) GetType() string {
	return string(c.Type)
}

// String returns a human-readable string (for debugging/logging)
func (c MessageContent) String() string {
	switch c.Type {
	case ContentTypeTextDelta:
		return c.Text
	case ContentTypeText:
		return c.Text
	case ContentTypeToolUse:
		args, _ := json.Marshal(c.Input)
		return fmt.Sprintf("%s:%s => %s", c.ID, c.Name, string(args))
	case ContentTypeToolResult:
		return fmt.Sprintf("Result[%s]: %s", c.ToolUseID, c.Content)
	default:
		return fmt.Sprintf("unknown content type: %s", c.Type)
	}
}

// Raw returns the entire struct as a generic interface
func (c MessageContent) Raw() interface{} {
	return c
}

// NewToolUseContent creates a tool use content message
func NewToolUseContent(id, name string, args json.RawMessage) *MessageContent {
	return &MessageContent{
		Type:  ContentTypeToolUse,
		ID:    id,
		Name:  name,
		Input: args,
	}
}

func NewToolResultContentInterface(toolUseID string, content interface{}) (*MessageContent, error) {
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("error marshalling tool result content: %w", err)
	}
	return &MessageContent{
		Type:      ContentTypeToolResult,
		ToolUseID: toolUseID,
		Content:   string(contentBytes),
	}, nil
}

// NewToolResultContent creates a tool result content message
func NewToolResultContent(toolUseID, content string) *MessageContent {
	return &MessageContent{
		Type:      ContentTypeToolResult,
		ToolUseID: toolUseID,
		Content:   content,
	}
}

// NewTextContent creates a text content message
func NewTextContent(text string) *MessageContent {
	return &MessageContent{
		Type: ContentTypeText,
		Text: text,
	}
}
