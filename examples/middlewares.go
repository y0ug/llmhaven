package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

// LoggingMiddleware creates a middleware that logs request and response details
func LoggingMiddleware() func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error) {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		// Log request
		reqDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Printf("Error dumping request: %v\n", err)
		} else {
			fmt.Printf("Request:\n%s\n", string(reqDump))
		}

		// Call the next middleware/handler
		resp, err := next(req)
		if err != nil {
			return resp, err
		}

		// Log response
		if resp != nil {
			// Read and restore the response body
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			respDump, err := httputil.DumpResponse(resp, false)
			if err != nil {
				fmt.Printf("Error dumping response: %v\n", err)
			} else {
				fmt.Printf("Response:\n%s\n", string(respDump))
				if resp.Request.Header.Get("Content-Type") == "application/json" {
					var obj map[string]interface{}
					err := json.Unmarshal([]byte(bodyBytes), &obj)
					if err == nil {
						bodyPretty, err := json.MarshalIndent(obj, "", "  ")
						if err == nil {
							bodyBytes = bodyPretty
						}
					}
				}
			}

			// Restore the response body again for subsequent readers
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		return resp, err
	}
}

func TimeitMiddleware() func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error) {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		start := time.Now()

		// Call the next middleware/handler
		resp, err := next(req)
		if err != nil {
			return resp, err
		}

		end := time.Now()

		fmt.Printf("Request took: %v\n", end.Sub(start))
		return resp, err
	}
}
