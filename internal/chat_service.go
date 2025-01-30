package internal

import (
	"context"

	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
)

type ChatService[Params any, Response any, Chunk any] interface {
	New(ctx context.Context, params Params, opts ...options.RequestOption) (Response, error)

	NewStreaming(
		ctx context.Context,
		params Params,
		opts ...options.RequestOption,
	) streaming.Streamer[Chunk]
}
