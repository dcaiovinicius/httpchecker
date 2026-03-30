package config

import (
	"errors"
	"flag"
	"time"
)

const (
	defaultConcurrency = 5
	defaultRequestTimeout = 10 * time.Second
	defaultGlobalTimeout  = 60 * time.Second
	defaultMaxRetries     = 2
	defaultUserAgent      = "Mozilla/5.0 (platform; rv:gecko-version) Gecko/gecko-trail Firefox/firefox-version"
	defaultDebug           = false
)


var (
	ErrInvalidConcurrency = errors.New("concurrency must be a positive integer")
	ErrInvalidRequestTimeout = errors.New("request timeout must be a positive duration")
	ErrInvalidGlobalTimeout  = errors.New("global timeout must be a positive duration")
	ErrInvalidMaxRetries     = errors.New("max retries must be a non-negative integer")
	ErrGlobalTimeoutLessThanRequestTimeout = errors.New("global timeout must be greater than or equal to request timeout")
)


type Config struct {
	Concurrency    int
	RequestTimeout time.Duration
	GlobalTimeout  time.Duration
	MaxRetries     int
	UserAgent      string
	Debug          bool
}

func (c *Config) Validate() error {
	if c.Concurrency <= 0 {
		return ErrInvalidConcurrency
	}

	if c.RequestTimeout <= 0 {
		return ErrInvalidRequestTimeout
	}
	
	if c.GlobalTimeout <= 0 {
		return ErrInvalidGlobalTimeout
	}

	if c.MaxRetries < 0 {
		return ErrInvalidMaxRetries
	}

	if c.GlobalTimeout < c.RequestTimeout {
		return ErrGlobalTimeoutLessThanRequestTimeout
	}
	
	return nil
}

func FromFlag() (*Config, error) {
	c := &Config{}

	flag.IntVar(&c.Concurrency, "concurrency", defaultConcurrency, "Number of concurrent workers")
	flag.DurationVar(&c.RequestTimeout, "request-timeout", defaultRequestTimeout, "Timeout for each HTTP request")
	flag.DurationVar(&c.GlobalTimeout, "global-timeout", defaultGlobalTimeout, "Overall timeout for the entire checking process")
	flag.IntVar(&c.MaxRetries, "max-retries", defaultMaxRetries, "Maximum number of retries for failed requests")
	flag.StringVar(&c.UserAgent, "user-agent", defaultUserAgent, "User-Agent string to use in HTTP requests")
	flag.BoolVar(&c.Debug, "debug", defaultDebug, "Enable debug logging")

	flag.Parse()

	if err := c.Validate(); err != nil {
		return nil, err
	}
	
	return c, nil
}
