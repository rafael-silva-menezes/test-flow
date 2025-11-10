package http

import (
	"context"
	stdhttp "net/http"
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

type Doer interface {
	Do(req *stdhttp.Request) (*stdhttp.Response, error)
}
