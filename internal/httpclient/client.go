package httpclient

import (
	"net/http"
	"time"
)

type DefaultClient struct {
	client     *http.Client
	retries    int
	retryDelay time.Duration
}

type Config struct {
	Timeout    time.Duration
	Retries    int
	RetryDelay time.Duration
}

func NewDefaultClient(config Config) *DefaultClient {
	return &DefaultClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		retries:    config.Retries,
		retryDelay: config.RetryDelay,
	}
}
