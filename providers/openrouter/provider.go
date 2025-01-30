package openrouter

import (
	"github.com/y0ug/llmhaven/chat"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/providers/openai"
)

type Provider struct {
	*openai.Provider
}

func New(opts ...options.RequestOption) chat.Provider {
	return &Provider{
		&openai.Provider{
			Client: NewClient(opts...).Client,
		},
	}
}
