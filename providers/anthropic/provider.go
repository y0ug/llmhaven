package anthropic

import (
	"context"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
)

type Provider struct {
	client *Client
}

func New(opts ...options.RequestOption) chat.Provider {
	return &Provider{
		client: NewClient(opts...),
	}
}

func (a *Provider) Send(
	ctx context.Context,
	params chat.ChatParams,
) (*chat.ChatResponse, error) {
	paramsProvider := BaseChatMessageNewParamsToAnthropic(params)
	am, err := a.client.Message.New(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}

	return AnthropicMessageToChatMessage(&am), nil
}

func (a *Provider) Stream(
	ctx context.Context,
	params chat.ChatParams,
) (streaming.Streamer[chat.EventStream], error) {
	paramsProvider := BaseChatMessageNewParamsToAnthropic(params)
	stream, err := a.client.Message.NewStreaming(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}
	return chat.NewProviderEventStream(
		stream,
		NewAnthropicEventHandler(),
	), nil
}

func (a *Provider) CountTokens(ctx context.Context,
	params chat.ChatParams,
) (int64, error) {
	paramsProvider := BaseChatMessageNewParamsToAnthropic(params)
	resp, err := a.client.Message.CountTokens(ctx, paramsProvider)
	if err != nil {
		return 0, err
	}
	return resp.InputTokens, nil
}
