package modelinfo

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock.go -package=modelinfo .  Provider

// A ModelInfoProvider implement the same method has CacheProvider[T]
// Expose the map of model name to model metadata in Val()
type Provider interface {
	Get(modelName string) (Getter, bool)
}

type Getter interface {
	Get(field string) (interface{}, bool)
	GetMaxTokens() int
	GetMaxInputTokens() int
	GetMaxOutputTokens() int
	GetInputCostPerToken() float64
	GetOutputCostPerToken() float64
	GetLiteLLMProvider() string
	GetCacheCreationInputTokenCost() *float64
	GetCacheReadInputTokenCost() *float64
	GetToolUseSystemPromptTokens() *int
}

var (
	_ Provider = (*LitellmModelInfoMap)(nil)
	_ Getter   = (*LitellmModelInfo)(nil)
)
