package config

import (
	"net/http"
	"testing"
	"time"
)

func TestRateLimitUpdate(t *testing.T) {
	rl := &RateLimit{}

	headers := http.Header{}
	headers.Set("X-Ratelimit-Limit-Requests", "10000")
	headers.Set("X-Ratelimit-Limit-Tokens", "200000")
	headers.Set("X-Ratelimit-Remaining-Requests", "9999")
	headers.Set("X-Ratelimit-Remaining-Tokens", "199976")
	headers.Set("X-Ratelimit-Reset-Requests", "8.64")
	headers.Set("X-Ratelimit-Reset-Tokens", "0.007")

	err := rl.Update(headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rl.LimitRequests != 10000 {
		t.Errorf("Expected LimitRequests 10000, got %d", rl.LimitRequests)
	}
	if rl.LimitTokens != 200000 {
		t.Errorf("Expected LimitTokens 200000, got %d", rl.LimitTokens)
	}
	if rl.RemainingRequests != 9999 {
		t.Errorf("Expected RemainingRequests 9999, got %d", rl.RemainingRequests)
	}
	if rl.RemainingTokens != 199976 {
		t.Errorf("Expected RemainingTokens 199976, got %d", rl.RemainingTokens)
	}
	if rl.ResetRequests != 8*time.Second+640*time.Millisecond {
		t.Errorf("Expected ResetRequests 8.64s, got %v", rl.ResetRequests)
	}
	if rl.ResetTokens != 7*time.Millisecond {
		t.Errorf("Expected ResetTokens 7ms, got %v", rl.ResetTokens)
	}
}
