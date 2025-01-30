package openai

import (
	"github.com/y0ug/llmhaven/http/config"
	"github.com/y0ug/llmhaven/http/options"
)

// WithOrganization returns a RequestOption that sets the client setting "organization".
func WithOrganization(value string) options.RequestOption {
	return func(r *config.RequestConfig) error {
		r.Organization = value
		return r.Apply(options.WithHeader("OpenAI-Organization", value))
	}
}

// WithProject returns a RequestOption that sets the client setting "project".
func WithProject(value string) options.RequestOption {
	return func(r *config.RequestConfig) error {
		r.Project = value
		return r.Apply(options.WithHeader("OpenAI-Project", value))
	}
}

// WithEnvironmentProductionOpenAI returns a RequestOption that sets the current
// environment to be the "production" environment. An environment specifies which base URL
// to use by default.
func WithEnvironmentProduction() options.RequestOption {
	return options.WithBaseURL("https://api.openai.com/v1/")
}
