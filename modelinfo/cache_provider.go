package modelinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock.go -package=modelinfo . CacheProvider

// CacheProvider defines the interface for cache operations.
// Expose a Val() that allow to recover the unmarshal json data
type CacheProvider[T any] interface {
	Val() *T
	Load(ctx context.Context) error
	Update(ctx context.Context) error
	Clear() error
	IsValid() bool
	Refresh(ctx context.Context) error
}

// NewCacheFileProvider creates a new CacheFileProvider instance.
func NewCacheFileProvider[T any](
	ctx context.Context,
	url, cacheFile string,
	cacheTTL time.Duration,
) (*CacheFileProvider[T], error) {
	c := &CacheFileProvider[T]{
		url:       url,
		cacheTTL:  cacheTTL,
		cacheFile: cacheFile,
	}
	if err := c.Load(ctx); err != nil {
		if err := c.Update(ctx); err != nil {
			return nil, fmt.Errorf("failed to load and update cache: %w", err)
		}
	}
	return c, nil
}

// CacheFileProvider implements the CacheProvider interface using a file as the cache.
type CacheFileProvider[T any] struct {
	url       string
	cacheTTL  time.Duration
	cacheFile string
	content   T
	mu        sync.RWMutex
}

// Get retrieves the cached content.
func (p *CacheFileProvider[T]) Val() *T {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return &p.content
}

// Load attempts to load the cache from the file. If the cache is invalid or absent, it does not automatically update.
func (p *CacheFileProvider[T]) Load(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cacheFile == "" {
		return fmt.Errorf("cache file not specified")
	}

	cacheDir := filepath.Dir(p.cacheFile)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if info, err := os.Stat(p.cacheFile); err == nil {
		if time.Since(info.ModTime()) < p.cacheTTL {
			data, err := os.ReadFile(p.cacheFile)
			if err != nil {
				return fmt.Errorf("failed to read cache file: %w", err)
			}

			var content T
			if err := json.Unmarshal(data, &content); err != nil {
				return fmt.Errorf("failed to parse downloaded data: %w", err)
			}
			p.content = content

			return nil
		}
	}

	return fmt.Errorf("cache is invalid or expired")
}

// Update fetches fresh data from the URL and updates the cache.
func (p *CacheFileProvider[T]) Update(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := download(ctx, p.url)
	if err != nil {
		return fmt.Errorf("failed to download data: %w", err)
	}

	var content T
	if err := json.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("failed to parse downloaded data: %w", err)
	}
	p.content = content

	if p.cacheFile != "" {
		if err := os.WriteFile(p.cacheFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write cache file: %w", err)
		}
	}

	return nil
}

// Clear invalidates the cache by resetting the content and removing the cache file.
func (p *CacheFileProvider[T]) Clear() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var content T
	p.content = content

	if p.cacheFile != "" {
		return os.Remove(p.cacheFile)
	}
	return nil
}

// IsValid checks whether the cache is still valid based on the TTL.
func (p *CacheFileProvider[T]) IsValid() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.cacheFile == "" {
		return false
	}

	if info, err := os.Stat(p.cacheFile); err == nil {
		return time.Since(info.ModTime()) < p.cacheTTL
	}

	return false
}

// Refresh updates the cache by fetching fresh data.
func (p *CacheFileProvider[T]) Refresh(ctx context.Context) error {
	return p.Update(ctx)
}

// download fetches data from the specified URL respecting the provided context.
func download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
