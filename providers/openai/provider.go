package openai

import (
	"context"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
)

type Provider struct {
	Client *Client
}

func New(opts ...options.RequestOption) chat.Provider {
	return &Provider{
		Client: NewClient(opts...),
	}
}

func (a *Provider) Send(
	ctx context.Context,
	params chat.ChatParams,
) (*chat.ChatResponse, error) {
	paramsProvider := ToChatCompletionNewParams(params)

	resp, err := a.Client.Chat.New(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}

	return ToChatResponse(&resp), nil
}

func (a *Provider) Stream(
	ctx context.Context,
	params chat.ChatParams,
) (streaming.Streamer[chat.EventStream], error) {
	paramsProvider := ToChatCompletionNewParams(params)

	stream, err := a.Client.Chat.NewStreaming(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}
	return chat.NewProviderEventStream(
		stream,
		NewOpenAIEventHandler(),
	), nil
}
