package openai

import (
	"os"

	"github.com/y0ug/llmhaven/http/client"
	"github.com/y0ug/llmhaven/http/options"
)

type Client struct {
	*client.BaseClient
	Chat *ChatCompletionService
}

func NewClient(opts ...options.RequestOption) (r *Client) {
	defaults := []options.RequestOption{
		WithEnvironmentProduction(),
	}
	if o, ok := os.LookupEnv("OPENAI_API_KEY"); ok {
		defaults = append(defaults, options.WithAuthToken(o))
	}
	if o, ok := os.LookupEnv("OPENAI_ORG_ID"); ok {
		defaults = append(defaults, WithOrganization(o))
	}
	if o, ok := os.LookupEnv("OPENAI_PROJECT_ID"); ok {
		defaults = append(defaults, WithProject(o))
	}
	r = &Client{
		BaseClient: client.NewBaseClient(NewAPIError, append(defaults, opts...)...),
	}

	r.Chat = NewChatCompletionService(r.Options...)

	return
}
