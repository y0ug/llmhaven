package deepseek

import (
	"context"

	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
	"github.com/y0ug/llmhaven/providers/openai"
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
	paramsProvider := openai.ToChatCompletionNewParams(params)

	resp, err := a.client.Chat.New(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}

	ret := &chat.ChatResponse{}
	ret.ID = resp.ID
	ret.Model = resp.Model
	ret.Usage = &chat.ChatUsage{}
	ret.Usage.InputTokens = resp.Usage.PromptTokens
	ret.Usage.OutputTokens = resp.Usage.CompletionTokens
	if len(resp.Choices) > 0 {
		for _, choice := range resp.Choices {
			c := chat.ChatChoice{}
			for _, call := range choice.Message.ToolCalls {
				c.Content = append(
					c.Content,
					openai.ToolCallToMessageContent(call),
				)
			}

			if choice.Message.Content != "" {
				c.Content = append(c.Content, chat.NewTextContent(choice.Message.Content))
			}

			// Role is not choice is our model
			c.Role = choice.Message.Role
			c.StopReason = openai.ToStopReason(choice.FinishReason)

			ret.Choice = append(ret.Choice, c)
		}
	}
	return ret, nil
}

func (a *Provider) Stream(
	ctx context.Context,
	params chat.ChatParams,
) (streaming.Streamer[chat.EventStream], error) {
	paramsProvider := openai.ToChatCompletionNewParams(params)

	stream, err := a.client.Chat.NewStreaming(ctx, paramsProvider)
	if err != nil {
		return nil, err
	}
	return chat.NewProviderEventStream(
		stream,
		openai.NewOpenAIEventHandler(),
	), nil
}
