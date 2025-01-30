package openrouter

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/y0ug/llmhaven/chat"
)

func TestSend(t *testing.T) {
	// Create a new provder with a mock client
	provder := New()
	ctx := context.Background()
	params := chat.ChatParams{
		Model:       "google/gemini-flash-1.5-8b",
		MaxTokens:   100,
		Temperature: 0.7,
		Messages: []*chat.ChatMessage{
			{
				Role: "user",
				Content: []*chat.MessageContent{
					chat.NewTextContent("Hello, how are you?"),
				},
			},
		},
	}

	// t.Skip("Skipping integration test - requires API key")

	response, err := provder.Send(ctx, params)

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, response)
	// Gemini don't set an response.ID
	// assert.NotEmpty(t, response.ID)
	assert.NotEmpty(t, response.Model)
	assert.Greater(t, response.Usage.InputTokens, 0)
	assert.Greater(t, response.Usage.OutputTokens, 0)

	fmt.Println(response.Choice[0].Content[0].String())
	fmt.Printf("Usage: %d %d\n", response.Usage.InputTokens, response.Usage.OutputTokens)
}
