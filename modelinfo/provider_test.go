package modelinfo

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestCacheFileProvider(t *testing.T) {
	t.Run("Test NewCacheFileProvider", func(t *testing.T) {
		provider, err := NewLitellm(context.Background(), "cache.json")
		if err != nil {
			log.Fatal(err)
		}

		models, ok := provider.Val().Get("gpt-4o")
		fmt.Println("models", models)
		if ok {
			if inputCostPerToken, ok := models.Get("input_cost_per_token"); ok {
				fmt.Println("input_cost_per_token", inputCostPerToken)
			}
			fmt.Println("MaxInputTokens", models.GetMaxInputTokens())
		}
		// fmt.Printf("provider %T %v\n")
		// fmt.Printf("provider %T %v\n", provider.Get(), provider)
		// Simple get
		// meta, exists := provider.Get()
		// if exists {
		// 	fmt.Printf("GPT-4o costs $%.10f per input token\n", meta.Get().InputCostPerToken)
		// }
		//

		// // Get with automatic refresh
		// meta, err = provider.GetModelWithFallback(context.Background(), "claude-3-opus")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Printf("Claude 3 Opus supports vision: %t\n", meta.SupportsVision)
	})
}
