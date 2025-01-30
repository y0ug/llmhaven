package anthropic

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/streaming"
	"go.uber.org/mock/gomock"
)

func skipIfNoAPIKey(t *testing.T) {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping integration test because ANTHROPIC_API_KEY is not set")
	}
}

func TestClientStreamIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewClient()
	ctx := context.Background()

	t.Run("ChatCompletion", func(t *testing.T) {
		params := MessageNewParams{
			Model:     "claude-3-5-sonnet-20241022",
			MaxTokens: 4096,
			Messages: []MessageParam{
				{
					Role: "user",
					Content: []*chat.MessageContent{chat.NewTextContent(
						"Write a 100 word essay on the topic of artificial intelligence",
					)},
				},
			},
			Temperature: 0,
		}
		stream, err := client.Message.NewStreaming(ctx, params)
		if err != nil {
			t.Fatalf("Failed to create chat completion stream: %v", err)
		}
		message := Message{}
		for stream.Next() {
			evt := stream.Current()
			message.Accumulate(evt)
			// evt := stream.Current()
			// switch evt := evt.(type) {
			// case *MessageStreamEvent:
			// 	message.Accumulate(*evt)
			// 	// fmt.Printf("%s ", evt.Type)
			// 	switch evt.Type {
			// 	case "message_start":
			// 	case "content_block_start":
			// 	case "content_block_delta":
			// 		// fmt.Printf("%v\n", evt.ContentBlock)
			// 		// fmt.Printf("Content: %v\n", evt.Delta)
			// 		fmt.Printf("%s", evt.Delta)
			// 	case "content_block_stop":
			// 	case "message_delta":
			// 		// fmt.Printf("%v\n", evt.ContentBlock)
			// 		// fmt.Printf("Content: %v\n", evt.Delta)
			// 	case "message_stop":
			// 	}
			// 	// fmt.Printf("\n")
			// default:
			// }
		}
		if stream.Err() != nil {
			fmt.Printf("Error: %v\n", stream.Err())
		}

		fmt.Printf("Message: %v\n", message.Content)
	})
}

func TestClientIntegration(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewClient()
	ctx := context.Background()

	t.Run("ChatCompletion", func(t *testing.T) {
		params := MessageNewParams{
			Model:     "claude-3-5-sonnet-20241022",
			MaxTokens: 4096,
			Messages: []MessageParam{
				{
					Role: "user",
					Content: []*chat.MessageContent{chat.NewTextContent(
						"Say hello in exactly 5 words",
					)},
				},
			},
			Temperature: 0,
		}

		message, err := client.Message.New(ctx, params)
		if err != nil {
			t.Fatalf("Failed to create chat completion: %v", err)
		}

		if len(message.Content) == 0 {
			t.Fatal("Expected at least one choice in response")
		}

		if message.Model == "" {
			t.Error("Expected model to be set in response")
		}

		if message.Usage.InputTokens == 0 {
			t.Error("Expected non-zero token usage")
		}

		t.Logf("Message: %v", message.Content[0])
	})

	t.Run("InvalidModel", func(t *testing.T) {
		params := MessageNewParams{
			Model:     "non-existent-model",
			MaxTokens: 4096,
			Messages: []MessageParam{
				{
					Role: "user",
					Content: []*chat.MessageContent{chat.NewTextContent(
						"Say hello in exactly 5 words",
					)},
				},
			},
		}

		_, err := client.Message.New(ctx, params)
		if err == nil {
			t.Error("Expected error for invalid model but got none")
		}
	})

	t.Run("StreamingOverloadedError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockStream := streaming.NewMockStreamer[chat.EventStream](mockCtrl)

		gomock.InOrder(
			// Mock the stream events sequence
			mockStream.EXPECT().Next().Return(true),
			mockStream.EXPECT().Current().Return(chat.EventStream{
				Type: "message_start",
				Message: &chat.ChatResponse{
					Choice: []chat.ChatChoice{
						{
							Content: []*chat.MessageContent{
								chat.NewTextContent("Looking at"),
							},
						},
					},
				},
			}),

			// Mock the error event
			mockStream.EXPECT().Next().Return(true),
			mockStream.EXPECT().Current().Return(chat.EventStream{
				Type:  "error",
				Delta: "Overloaded",
			}),

			// Final sequence after error
			// mockStream.EXPECT().Next().Return(false),
			// mockStream.EXPECT().Err().Return(fmt.Errorf("Overloaded")),
			// mockStream.EXPECT().Close().Return(nil),
		)

		// Create a mock client that returns our mock stream
		mockClient := chat.NewMockProvider(mockCtrl)
		mockClient.EXPECT().
			Stream(gomock.Any(), gomock.Any()).
			Return(mockStream, nil)

		params := chat.NewChatParams(
			chat.WithModel("claude-3-5-sonnet-20241022"),
			chat.WithMaxTokens(4096),
			chat.WithMessages(
				chat.NewUserMessage("Test message"),
			),
		)

		_, err := HandleLLMConversation(context.Background(), mockClient, *params)
		if err == nil {
			t.Error("Expected overloaded error but got none")
		}

		if err != nil && err.Error() != "Overloaded" {
			t.Errorf("Expected 'Overloaded' error but got: %v", err)
		}
	})
}

func HandleLLMConversation(
	ctx context.Context,
	provider chat.Provider,
	params chat.ChatParams,
) (*chat.ChatResponse, error) {
	var msg *chat.ChatResponse
	for {

		stream, err := provider.Stream(ctx, params)
		if err != nil {
			log.Printf("Error streaming: %v", err)
			return nil, err
		}

		eventCh := make(chan chat.EventStream)

		// llmclient.ConsumeStreamIO(ctx, stream, os.Stdout)
		go func() {
			// llmclient.ConsumeStreamIO(ctx, stream, os.Stdout)
			if err := chat.StreamChatMessageToChannel(ctx, stream, eventCh); err != nil {
				if err != context.Canceled {
					log.Printf("Error consuming stream: %v", err)
					close(eventCh)
					return
				}
			}
		}()

		msg, err = processStream(ctx, os.Stdout, eventCh)
		log.Printf("msg: %v", msg)
		if err != nil {
			log.Printf("Error processing stream: %v", err)
			return nil, err
		}

		if msg == nil {
			log.Printf("No message returned")
			return nil, nil
		}
		fmt.Printf("\nUsage: %d %d\n", msg.Usage.InputTokens, msg.Usage.OutputTokens)

		params.Messages = append(params.Messages, msg.ToMessageParams())
		toolResults := make([]*chat.MessageContent, 0)
		// for _, choice := range msg.Choice {
		choice := msg.Choice[0]
		for _, content := range choice.Content {
			if content.Type == "tool_use" {
				log.Printf(
					"%s execution: %s with \"%s\"",
					content.ID,
					content.Name,
					string(content.Input),
				)
				switch content.Name {
				default:
					log.Printf("Unknown tool: %s", content.Name)
				}

			}
		}
		// }
		if len(toolResults) == 0 {
			break
		}

		// if params.N != nil {
		// 	*params.N = 1
		// }

		params.Messages = append(params.Messages, chat.NewMessage("user", toolResults...))
	}
	return msg, nil
}

func processStream(
	ctx context.Context,
	w io.Writer,
	ch <-chan chat.EventStream,
) (*chat.ChatResponse, error) {
	var cm *chat.ChatResponse
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case set, ok := <-ch:
			if !ok {
				return cm, nil
			}
			switch set.Type {
			case "text_delta":
				fmt.Fprintf(w, "%v", set.Delta)
			case "message_stop":
				cm = set.Message
			case "error":
				cm = set.Message
				return cm, fmt.Errorf("%v", set.Delta)
			}
		}
	}
}
