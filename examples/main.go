package examples

import (
	"context"
	"fmt"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
)

func main() {
	ctx := context.Background()
	model := ""
	provider := ""

	messages := []*chat.ChatMessage{
		chat.NewSystemMessage("You're a weather expert."),
		chat.NewUserMessage("What's the weather like in Paris?"),
	}

	params := chat.NewChatParams(
		chat.WithModel(model),
		chat.WithMessages(messages...),
		// chat.WithTools(tools...))
	)

	// Request options for the API, for example we can dump API request/response
	opts := []options.RequestOption{
		options.WithMiddleware(LoggingMiddleware()),
		options.WithMiddleware(TimeitMiddleware()),
	}

	resp, err := ChatCompletion(ctx, provider, *params, opts)
	if err != nil {
		fmt.Printf("Failed to get response: %v", err)
	}

	if resp.HasContent() {
		fmt.Printf("Response: %s", resp.Choice[0].Content[0].Content)
	}
}
