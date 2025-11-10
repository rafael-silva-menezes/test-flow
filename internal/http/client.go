package http

import (
	"bytes"
	"context"
	"fmt"
	stdhttp "net/http"
	"time"
)

type DefaultClient struct {
	client *stdhttp.Client
	config Config
}
type Config struct {
	DefaultTimeout time.Duration
	Retries        int
	RetryDelay     time.Duration
}

func NewDefaultClient(config Config) *DefaultClient {
	return &DefaultClient{
		client: &stdhttp.Client{
			Timeout: config.DefaultTimeout,
		},
		config: config,
	}
}

func (c *DefaultClient) Get(ctx context.Context, req Request) (Response, error) {
	return c.doRequest(ctx, req, stdhttp.MethodGet)
}

func (c *DefaultClient) Post(ctx context.Context, req Request) (Response, error) {
	return c.doRequest(ctx, req, stdhttp.MethodPost)
}

func (c *DefaultClient) doRequest(ctx context.Context, req Request, method string) (Response, error) {
	timeout := c.config.DefaultTimeout
	if req.Timeout > 0 {
		timeout = req.Timeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url, body := req.URL, bytes.NewReader(req.Body)
	httpReq, err := stdhttp.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		resp, err := c.client.Do(httpReq)
		if err != nil {
			lastErr = err
			time.Sleep(c.config.RetryDelay)
			continue
		}

		defer resp.Body.Close()
	}
}
