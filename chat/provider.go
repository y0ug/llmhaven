package chat

import (
	"context"

	"github.com/y0ug/llmhaven/http/streaming"
)

// ChatProvider
//
//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock.go -package=chat .  Provider
type Provider interface {
	// For a single-turn request
	Send(ctx context.Context, params ChatParams) (*ChatResponse, error)

	// For streaming support
	Stream(
		ctx context.Context,
		params ChatParams,
	) (streaming.Streamer[EventStream], error)
}
