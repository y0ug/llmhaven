package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/y0ug/llmhaven/http/config"
	"github.com/y0ug/llmhaven/http/options"
	"github.com/y0ug/llmhaven/http/streaming"
)

// GenericChatService is a generic base implementation of ChatService.
type GenericChatService[Params any, Response any, Chunk any] struct {
	Options  []options.RequestOption
	NewError config.NewAPIError
	Endpoint string
}

// New creates a new chat completion.
func (svc *GenericChatService[Params, Response, Chunk]) New(
	ctx context.Context,
	params Params,
	opts ...options.RequestOption,
) (Response, error) {
	var res Response
	combinedOpts := append(svc.Options, opts...)
	path := svc.Endpoint

	err := config.ExecuteNewRequest(
		ctx,
		http.MethodPost,
		path,
		params,
		&res,
		svc.NewError,
		combinedOpts...,
	)
	return res, err
}

// NewStreaming creates a new streaming chat completion.
func (svc *GenericChatService[Params, Response, Chunk]) NewStreaming(
	ctx context.Context,
	params Params,
	opts ...options.RequestOption,
) (streaming.Streamer[Chunk], error) {
	combinedOpts := append(svc.Options, opts...)
	combinedOpts = append(
		[]options.RequestOption{
			options.WithJSONSet("stream", true),
			options.WithJSONSet("stream_options", struct {
				IncludeUsage bool `json:"include_usage"`
			}{IncludeUsage: true}),
		},
		combinedOpts...)
	path := svc.Endpoint

	var raw *http.Response
	err := config.ExecuteNewRequest(
		ctx,
		http.MethodPost,
		path,
		params,
		&raw,
		svc.NewError,
		combinedOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing new request streaming: %w", err)
	}
	return streaming.NewStream(
		streaming.NewDecoderSSE(raw),
		streaming.NewGenericStreamHandler[Chunk](),
	), nil
}
