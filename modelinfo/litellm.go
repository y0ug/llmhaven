package modelinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

/*
models, ok := provider.Val().Get("gpt-4o")
fmt.Println("models", models)

	if ok {
		if inputCostPerToken, ok := models.Get("input_cost_per_token"); ok {
			fmt.Println("input_cost_per_token", inputCostPerToken)
		}
		fmt.Println("MaxInputTokens", models.MaxInputTokens)
	}
*/
var littellmDataURL = "https://raw.githubusercontent.com/BerriAI/litellm/main/model_prices_and_context_window.json"

// It's used the CacheFileProvider to expose the data under Val()
//
//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock.go -package=modelinfo .  LitellmFileInfoCacheProvider
type LitellmCacheProvider struct {
	*CacheFileProvider[LitellmModelInfoMap]
}

func NewLitellm(
	ctx context.Context,
	cacheFile string,
) (*LitellmCacheProvider, error) {
	cfp, err := NewCacheFileProvider[LitellmModelInfoMap](
		ctx,
		littellmDataURL,
		cacheFile,
		24*time.Hour,
	)
	if err != nil {
		return nil, err
	}
	return &LitellmCacheProvider{cfp}, nil
}

// LitellmModelInfoFull represents the complete Litellm model metadata.
type LitellmModelInfoFull struct {
	MaxTokens          int     `json:"max_tokens"`
	MaxInputTokens     int     `json:"max_input_tokens"`
	MaxOutputTokens    int     `json:"max_output_tokens"`
	InputCostPerToken  float64 `json:"input_cost_per_token"`
	OutputCostPerToken float64 `json:"output_cost_per_token"`
	LiteLLMProvider    string  `json:"litellm_provider"`
	Mode               string  `json:"mode"`

	// Optional fields
	CacheCreationInputTokenCost *float64 `json:"cache_creation_input_token_cost,omitempty"`
	CacheReadInputTokenCost     *float64 `json:"cache_read_input_token_cost,omitempty"`
	ToolUseSystemPromptTokens   *int     `json:"tool_use_system_prompt_tokens,omitempty"`

	// Optional feature flags
	SupportsFunctionCalling         bool `json:"supports_function_calling,omitempty"`
	SupportsParallelFunctionCalling bool `json:"supports_parallel_function_calling,omitempty"`
	SupportsVision                  bool `json:"supports_vision,omitempty"`
	SupportsAudioInput              bool `json:"supports_audio_input,omitempty"`
	SupportsAudioOutput             bool `json:"supports_audio_output,omitempty"`
	SupportsPromptCaching           bool `json:"supports_prompt_caching,omitempty"`
	SupportsResponseSchema          bool `json:"supports_response_schema,omitempty"`
	SupportsSystemMessages          bool `json:"supports_system_messages,omitempty"`
	SupportsAssistantPrefill        bool `json:"supports_assistant_prefill,omitempty"`
}

// Represents the metadata for a model provide by littlellm model info json
type LitellmModelInfo LitellmModelInfoFull

// They provide map of model name to model metadata
type LitellmModelInfoMap map[string]LitellmModelInfo

func (v *LitellmModelInfo) GetMaxTokens() int {
	return v.MaxTokens
}

func (v *LitellmModelInfo) GetMaxInputTokens() int {
	return v.MaxInputTokens
}

func (v *LitellmModelInfo) GetMaxOutputTokens() int {
	return v.MaxOutputTokens
}

func (v *LitellmModelInfo) GetInputCostPerToken() float64 {
	return v.InputCostPerToken
}

func (v *LitellmModelInfo) GetOutputCostPerToken() float64 {
	return v.OutputCostPerToken
}

func (v *LitellmModelInfo) GetLiteLLMProvider() string {
	return v.LiteLLMProvider
}

func (v *LitellmModelInfo) GetCacheCreationInputTokenCost() *float64 {
	return v.CacheCreationInputTokenCost
}

func (v *LitellmModelInfo) GetCacheReadInputTokenCost() *float64 {
	return v.CacheReadInputTokenCost
}

func (v *LitellmModelInfo) GetToolUseSystemPromptTokens() *int {
	return v.ToolUseSystemPromptTokens
}

// Reflect the field name to the json name
func (mc LitellmModelInfo) Get(name string) (interface{}, bool) {
	v := reflect.ValueOf(mc).Elem()
	fieldVal := v.FieldByNameFunc(func(s string) bool {
		field, ok := reflect.TypeOf(mc).Elem().FieldByName(s)
		return ok && field.Tag.Get("json") == name
	})

	if !fieldVal.IsValid() {
		return nil, false
	}

	return fieldVal.Interface(), true
}

// Return the metadata for a model by name
func (p *LitellmModelInfoMap) Get(modelName string) (Getter, bool) {
	val, ok := (*p)[modelName]
	return &val, ok
}

// UnmarshalJSON implements the json.Unmarshaler interface for LitellmMetadataMap
// Ignore the key sample_spec
// Ignore entry when the unmarshal fails
func (p *LitellmModelInfoMap) UnmarshalJSON(b []byte) error {
	if *p == nil {
		*p = make(LitellmModelInfoMap, 0)
	}
	// First, unmarshal into a map of raw messages
	temp := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal into temporary map: %w", err)
	}

	for key, raw := range temp {
		if key == "sample_spec" {
			continue
		}
		var tmp LitellmModelInfo
		if err := json.Unmarshal(raw, &tmp); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to parse info for model %s: %v\n", key, err)
			continue
		}
		(*p)[key] = tmp
	}

	return nil
}
