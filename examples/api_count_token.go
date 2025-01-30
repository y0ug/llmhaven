package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/y0ug/llmhaven"
	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/providers/anthropic"
)

func ApiCountToken(
	ctx context.Context,
	provider string,
	params chat.ChatParams,
	opts []options.RequestOption,
) (int64, error) {
	ctxRequest, cancelFn := context.WithTimeout(ctx, 10*time.Second)
	defer cancelFn()

	llm, err := llmhaven.New(provider, opts...)
	if err != nil {
		return 0, fmt.Errorf("Failed to create provider: %v", err)
	}

	if llm, ok := llm.(*anthropic.Provider); ok {
		tokens, err := llm.CountTokens(ctxRequest, params)
		if err != nil {
			return 0, fmt.Errorf("Failed to count tokens: %v", err)
		}
		return tokens, nil
	}

	return 0, fmt.Errorf("Provider not supported")
}
