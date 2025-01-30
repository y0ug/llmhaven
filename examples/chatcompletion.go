package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/y0ug/llmhaven"
	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
)

func ChatCompletion(
	ctx context.Context,
	provider string,
	params chat.ChatParams,
	opts []options.RequestOption,
) (*chat.ChatResponse, error) {
	ctxRequest, cancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFn()

	llm, err := llmhaven.New(provider, opts...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create provider: %v", err)
	}

	resp, err := llm.Send(ctxRequest, params)
	if err != nil {
		return nil, fmt.Errorf("Failed to send message: %v", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("Response is nil")
	}

	return resp, nil
}
