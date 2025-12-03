package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDefaultFetcher_Fetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "test-ua/1.0" {
			t.Errorf("Expected User-Agent test-ua/1.0, got %q", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><a href="/link1">Link1</a></body></html>`))
	}))
	defer server.Close()

	fetcher := NewDefaultFetcher("test-ua/1.0", 10*time.Second)
	ctx := context.Background()

	links, statusCode, err := fetcher.Fetch(ctx, server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if statusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", statusCode)
	}
	if len(links) == 0 {
		t.Error("Expected some links extracted")
	}
}

func TestDefaultFetcher_Fetch_InvalidURL(t *testing.T) {
	fetcher := NewDefaultFetcher("test-ua/1.0", 10*time.Second)
	ctx := context.Background()

	_, _, err := fetcher.Fetch(ctx, "invalid-url://test")

	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDefaultFetcher_Fetch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	fetcher := NewDefaultFetcher("test-ua/1.0", 10*time.Second)
	ctx := context.Background()

	links, statusCode, err := fetcher.Fetch(ctx, server.URL)

	if err != nil {
		t.Errorf("Expected no error on HTTP error, got %v", err)
	}
	if statusCode != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", statusCode)
	}
	if len(links) != 0 {
		t.Error("Expected empty links on HTTP error")
	}
}

func TestDefaultFetcher_Fetch_ContextTimeout(t *testing.T) {
	fetcher := NewDefaultFetcher("test-ua/1.0", 10*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	_, _, err := fetcher.Fetch(ctx, "https://httpbin.org/delay/1")

	if err == nil {
		t.Error("Expected context timeout error")
	}
}
