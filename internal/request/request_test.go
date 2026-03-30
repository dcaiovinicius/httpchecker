package request

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"urlchecker/internal/config"
)

func TestShouldRetry(t *testing.T) {
	t.Parallel()
	if shouldRetry(200, nil) {
		t.Fatal("success should not retry")
	}
	if !shouldRetry(0, fmt.Errorf("network")) {
		t.Fatal("network error should retry")
	}
	if !shouldRetry(http.StatusInternalServerError, fmt.Errorf("HTTP 500")) {
		t.Fatal("5xx should retry")
	}
	if !shouldRetry(http.StatusTooManyRequests, fmt.Errorf("HTTP 429")) {
		t.Fatal("429 should retry")
	}
	if shouldRetry(http.StatusNotFound, fmt.Errorf("HTTP 404")) {
		t.Fatal("404 should not retry")
	}
}

func TestChecker_Check_OK(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	cfg := testConfig(0)
	c := NewChecker(cfg)
	if err := c.Check(context.Background(), srv.URL); err != nil {
		t.Fatal(err)
	}
}

func TestChecker_Check_NotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	cfg := testConfig(2)
	c := NewChecker(cfg)
	if err := c.Check(context.Background(), srv.URL); err == nil {
		t.Fatal("expected error for 404")
	}
}

func testConfig(maxRetries int) *config.Config {
	return &config.Config{
		Concurrency:    2,
		RequestTimeout: 5 * time.Second,
		GlobalTimeout:  30 * time.Second,
		MaxRetries:     maxRetries,
		UserAgent:      "urlchecker/test",
		Debug:          false,
	}
}
