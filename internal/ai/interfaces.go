package ai

import (
	"context"
	"net/http"
)

type AIClient interface {
	GenerateTest(ctx context.Context, prompt string) (AIResponse, error)
}

type AIParser interface {
	ParseAIResponse(raw string) (AIResponse, error)
}

type AIResponse struct {
	TestName string `json:"test_name"`
	Code     string `json:"code"`
}

// AIError represents an error returned by the AI service.
type AIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *AIError) Error() string {
	return e.Message
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
