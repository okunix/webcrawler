package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

var (
	userAgentHeader = http.CanonicalHeaderKey("User-Agent")
)

type Fetcher interface {
	Fetch(ctx context.Context, url string) ([]string, int, error)
}

type DefaultFetcher struct {
	userAgent string
	timeout   time.Duration
}

func NewDefaultFetcher(userAgent string, timeout time.Duration) Fetcher {
	return &DefaultFetcher{userAgent: userAgent, timeout: timeout}
}

func (f *DefaultFetcher) Fetch(ctx context.Context, url string) ([]string, int, error) {
	timeout, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Add(userAgentHeader, f.userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	baseURL, _ := BaseURL(url)
	links := ExtractLinks(string(body), baseURL)
	return links, resp.StatusCode, nil
}
