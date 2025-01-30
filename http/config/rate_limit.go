package config

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimit holds the rate limit information extracted from response headers.
type RateLimit struct {
	LimitRequests     int
	LimitTokens       int
	RemainingRequests int
	RemainingTokens   int
	ResetRequests     time.Duration
	ResetTokens       time.Duration
	mu                sync.RWMutex
}

// Update updates the RateLimit fields based on the provided headers.
func (rl *RateLimit) Update(headers http.Header) error {
	limitRequests, err := strconv.Atoi(headers.Get("X-Ratelimit-Limit-Requests"))
	if err != nil {
		return err
	}

	limitTokens, err := strconv.Atoi(headers.Get("X-Ratelimit-Limit-Tokens"))
	if err != nil {
		return err
	}

	remainingRequests, err := strconv.Atoi(headers.Get("X-Ratelimit-Remaining-Requests"))
	if err != nil {
		return err
	}

	remainingTokens, err := strconv.Atoi(headers.Get("X-Ratelimit-Remaining-Tokens"))
	if err != nil {
		return err
	}

	resetRequestsSec, err := strconv.ParseFloat(headers.Get("X-Ratelimit-Reset-Requests"), 64)
	if err != nil {
		return err
	}

	resetTokensSec, err := strconv.ParseFloat(headers.Get("X-Ratelimit-Reset-Tokens"), 64)
	if err != nil {
		return err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.LimitRequests = limitRequests
	rl.LimitTokens = limitTokens
	rl.RemainingRequests = remainingRequests
	rl.RemainingTokens = remainingTokens
	rl.ResetRequests = time.Duration(resetRequestsSec * float64(time.Second))
	rl.ResetTokens = time.Duration(resetTokensSec * float64(time.Second))

	return nil
}

// GetSnapshot returns a copy of the current RateLimit status.
func (rl *RateLimit) GetSnapshot() RateLimit {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return RateLimit{
		LimitRequests:     rl.LimitRequests,
		LimitTokens:       rl.LimitTokens,
		RemainingRequests: rl.RemainingRequests,
		RemainingTokens:   rl.RemainingTokens,
		ResetRequests:     rl.ResetRequests,
		ResetTokens:       rl.ResetTokens,
	}
}
