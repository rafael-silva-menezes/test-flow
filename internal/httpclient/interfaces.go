package httpclient

import (
	"context"
	"time"
)

type Request struct {
	URL     string
	Body    []byte
	Headers map[string]string
	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

type HTTPClient interface {
	Get(ctx context.Context, req Request) (Response, error)
	Post(ctx context.Context, req Request) (Response, error)
}
