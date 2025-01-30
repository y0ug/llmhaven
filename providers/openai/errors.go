package openai

import (
	"encoding/json"
	"net/http"

	"github.com/y0ug/llmhaven/http/errors"
)

type APIError struct {
	errors.APIErrorBase
	Code    string `json:"code,required,nullable"`
	Message string `json:"message,required"`
	Param   string `json:"param,required,nullable"`
	Type    string `json:"type,required"`
	JSON    string `json:"-"`
}

func (r *APIError) UnmarshalJSON(data []byte) (err error) {
	r.JSON = string(data)
	type Alias APIError
	return json.Unmarshal(data, (*Alias)(r))
}

func NewAPIError(resp *http.Response, req *http.Request) errors.APIError {
	return &APIError{
		APIErrorBase: errors.APIErrorBase{
			StatusCode: resp.StatusCode,
			Request:    req,
			Response:   resp,
		},
	}
}
