package modelinfo

import (
	"context"
	"fmt"
	"strings"
)

type Model struct {
	Provider string
	Name     string
	info     Getter
}

func New(ctx context.Context, cacheFile string) (Provider, error) {
	cachedProvider, err := NewLitellm(ctx, cacheFile)
	if err != nil {
		return nil, err
	}
	return cachedProvider.Val(), nil
}

// String returns the string representation of the model
func (m *Model) String() string {
	return fmt.Sprintf("%s/%s", m.Provider, m.Name)
}

func (m *Model) Info() Getter {
	return m.info
}

// inferProvider attempts to determine the provider based on model name patterns
func inferProvider(modelName string) string {
	modelName = strings.ToLower(modelName)

	switch {
	case strings.HasPrefix(modelName, "claude"):
		return "anthropic"
	case strings.HasPrefix(modelName, "deepseek"):
		return "deepseek"
	case strings.HasPrefix(modelName, "gpt"):
		return "openai"
	case strings.HasPrefix(modelName, "gemini"):
		return "google"
	case strings.HasPrefix(modelName, "mistral"):
		return "mistral"
	case strings.HasPrefix(modelName, "llama"):
		return "meta"
	default:
		return ""
	}
}

// Get parses a model string in the format "provider/model"
func Get(modelStr string, infoProviders Provider) (*Model, error) {
	if modelStr == "" {
		return nil, fmt.Errorf("empty model string")
	}

	parts := strings.Split(modelStr, "/")

	if len(parts) < 2 {
		// If model not found in info or no providers, try to infer provider from model name
		provider := inferProvider(modelStr)

		// For models without explicit provider prefix, try to get info if providers available
		if infoProviders != nil {
			info, ok := infoProviders.Get(modelStr)
			if !ok {
				if info.GetLiteLLMProvider() != "" {
					provider = info.GetLiteLLMProvider()
				}
				if provider == "" {
					return nil, fmt.Errorf("could not determine provider for model: %s", modelStr)
				}
				return &Model{
					Provider: provider,
					Name:     modelStr,
					info:     info,
				}, nil
			}
		}
		if provider != "" {
			return &Model{
				Provider: provider,
				Name:     modelStr,
				info:     nil,
			}, nil
		}
		return nil, fmt.Errorf("could not determine provider for model: %s", modelStr)
	}

	provider := parts[0]
	name := strings.Join(parts[1:], "/")

	var result Getter
	if infoProviders != nil {
		info, ok := infoProviders.Get(modelStr)
		// Even if we failed to load the info we can still return the model
		// without info
		if ok {
			result = info
		}
	}
	return &Model{
		Provider: provider,
		Name:     name,
		info:     result,
	}, nil
}
