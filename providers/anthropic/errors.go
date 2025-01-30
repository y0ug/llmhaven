package anthropic

import (
	"encoding/json"
	"net/http"

	"github.com/y0ug/llmhaven/http/errors"
)

func NewError(resp *http.Response, req *http.Request) errors.APIError {
	return &APIError{
		APIErrorBase: errors.APIErrorBase{
			StatusCode: resp.StatusCode,
			Request:    req,
			Response:   resp,
		},
	}
}

type APIError struct {
	errors.APIErrorBase
	ExtraFields map[string]interface{} `json:"-"`
}

func (r *APIError) UnmarshalJSON(data []byte) (err error) {
	r.JSON = string(data)
	r.ExtraFields = make(map[string]interface{})
	return json.Unmarshal(data, &r.ExtraFields)
}
