package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"dcaiovinicius/httpchecker/internal/config"
	"dcaiovinicius/httpchecker/internal/logger"
)

type Checker struct {
	client *http.Client
	config *config.Config
}

func NewChecker(cfg *config.Config) *Checker {
	return &Checker{
		client: &http.Client{Timeout: cfg.RequestTimeout},
		config: cfg,
	}
}

func (c *Checker) Check(ctx context.Context, rawURL string) error {
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if attempt > 0 {
			logger.Info("retry %d/%d for %s", attempt, c.config.MaxRetries, rawURL)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}

		code, err := c.checkOnce(ctx, rawURL)
		if err == nil {
			return nil
		}
		lastErr = err
		if !shouldRetry(code, err) {
			return err
		}
	}
	return lastErr
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 150 * time.Millisecond
	if d > 2*time.Second {
		return 2 * time.Second
	}
	return d
}

func shouldRetry(statusCode int, err error) bool {
	if err == nil {
		return false
	}
	if statusCode == http.StatusTooManyRequests {
		return true
	}
	if statusCode >= 500 {
		return true
	}
	return statusCode == 0
}

func statusType(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "SUCCESS"
	case code >= 300 && code < 400:
		return "REDIRECT"
	case code >= 400 && code < 500:
		return "CLIENT_ERROR"
	case code >= 500:
		return "SERVER_ERROR"
	default:
		return "UNKNOWN"
	}
}

func (c *Checker) checkOnce(ctx context.Context, rawURL string) (statusCode int, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", c.config.UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
		return resp.StatusCode, fmt.Errorf("read body: %w", copyErr)
	}

	code := resp.StatusCode
	logger.Notice("[%s] %s -> %d %s", statusType(code), rawURL, code, http.StatusText(code))

	if code >= 200 && code < 400 {
		return code, nil
	}

	if code == http.StatusTooManyRequests {
		return code, fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
	}
	if code >= 500 {
		return code, fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
	}
	if code >= 400 {
		return code, fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
	}
	return code, fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
}
