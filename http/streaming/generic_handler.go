package streaming

import (
	"encoding/json"
	"fmt"
)

// GenericStreamHandler implements StreamHandler for basic streaming responses
type GenericStreamHandler[T any, TypeIn Event] struct{}

func NewGenericStreamHandler[T any, TypeIn Event]() *GenericStreamHandler[T, Event] {
	return &GenericStreamHandler[T, Event]{}
}

func (h *GenericStreamHandler[T, TypeIn]) HandleEvent(event Event) (T, error) {
	var result T

	if len(event.Data) == 0 {
		return result, nil
	}

	if err := json.Unmarshal(event.Data, &result); err != nil {
		return result, fmt.Errorf("error unmarshalling event: %w %s", err, event.Data)
	}

	return result, nil
}

func (h *GenericStreamHandler[T, TypeIn]) ShouldContinue(event Event) bool {
	// fmt.Printf("[DONE] != '%s'\n", string(event.Data))
	return string(event.Data) != "[DONE]"
}
