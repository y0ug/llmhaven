package gemini

import (
	"os"

	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/providers/openai"
)

type Client struct {
	*openai.Client
}

func WithEnvironmentProduction() options.RequestOption {
	return options.WithBaseURL("https://generativelanguage.googleapis.com/v1beta/openai/")
}

func NewClient(opts ...options.RequestOption) *Client {
	defaults := []options.RequestOption{
		WithEnvironmentProduction(),
	}
	if o, ok := os.LookupEnv("GEMINI_API_KEY"); ok {
		defaults = append(defaults, options.WithAuthToken(o))
	}
	opts = append(defaults, opts...)
	r := &Client{
		Client: openai.NewClient(opts...),
	}

	return r
}
