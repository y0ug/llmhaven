package deepseek

import (
	"encoding/json"
	"os"

	"github.com/y0ug/llmhaven/http/client"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/internal"
	"github.com/y0ug/llmhaven/providers/openai"
)

type DeepSeekUsage struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	TotalTokens           int `json:"total_tokens"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}

type Client struct {
	*client.BaseClient
	Chat *ChatCompletionService
}

func WithEnvironmentProduction() options.RequestOption {
	return options.WithBaseURL("https://api.deepseek.com/v1/")
}

func NewClient(opts ...options.RequestOption) (r *Client) {
	defaults := []options.RequestOption{
		WithEnvironmentProduction(),
		// requestoption.WithMiddleware(middleware.LoggingMiddleware()),
	}
	if o, ok := os.LookupEnv("DEEPSEEK_API_KEY"); ok {
		defaults = append(defaults, options.WithAuthToken(o))
	}
	opts = append(defaults, opts...)
	r = &Client{
		BaseClient: client.NewBaseClient(openai.NewAPIError, opts...),
	}

	r.Chat = NewChatCompletionService(r.BaseClient.Options...)
	return r
}

// We define the new ChatCompletionService that embeds the BaseChatService
type ChatCompletionService struct {
	*internal.GenericChatService[openai.ChatCompletionNewParams, ChatCompletion, openai.ChatCompletionChunk]
}

// Create our custom service that reuses openai.NewChatCompletionService under the hood
func NewChatCompletionService(opts ...options.RequestOption) *ChatCompletionService {
	baseService := &internal.GenericChatService[openai.ChatCompletionNewParams, ChatCompletion, openai.ChatCompletionChunk]{
		Options:  opts,
		NewError: openai.NewAPIError,
		Endpoint: "chat/completions",
	}

	return &ChatCompletionService{
		GenericChatService: baseService,
	}
}

type ChatCompletion struct {
	// Embed all the OpenAI fields.
	openai.ChatCompletion
	// Override usage with your DeepSeek usage struct
	Usage DeepSeekUsage `json:"usage"`
}

// Custom Unmarshal to fill our DeepSeekUsage, while still reusing the embedded fields.
func (r *ChatCompletion) UnmarshalJSON(data []byte) error {
	type Alias ChatCompletion
	if err := json.Unmarshal(data, (*Alias)(r)); err != nil {
		return err
	}

	// TODO: due to golang embedded overide
	// we have to override the usage field
	// Then, parse out usage with a simple JSON map
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Now unmarshal 'usage' into your custom type
	if usageBytes, ok := raw["usage"]; ok {
		if err := json.Unmarshal(usageBytes, &r.Usage); err != nil {
			return err
		}
	}
	return nil
}
