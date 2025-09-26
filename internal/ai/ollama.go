package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	model   string
	baseURL string
	stream  bool
	parser  AIParser
	client  HTTPDoer
}

func NewOllamaClient(model, baseURL string, timeout time.Duration, stream bool, parser AIParser) *OllamaClient {
	return &OllamaClient{
		model:   model,
		baseURL: baseURL,
		stream:  stream,
		parser:  parser,
		client:  &http.Client{Timeout: timeout},
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (o *OllamaClient) generateRaw(ctx context.Context, prompt string) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", &AIError{
			Message:    "prompt cannot be empty",
			StatusCode: http.StatusBadRequest,
		}
	}

	body, err := o.buildPayload(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to build payload: %w", err)
	}

	req, err := o.buildRequest(ctx, body)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return o.handleResponse(resp)
}

func (o *OllamaClient) buildPayload(prompt string) ([]byte, error) {
	payload := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: o.stream,
	}
	return json.Marshal(payload)
}

func (o *OllamaClient) buildRequest(ctx context.Context, body []byte) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/generate", o.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (o *OllamaClient) handleResponse(resp *http.Response) (string, error) {
	const maxBody = 5 * 1024 * 1024 // 5 MB
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxBody))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Truncate error message if too large
		msg := string(respBody)
		if len(msg) > 2048 {
			msg = msg[:2048] + "... [truncated]"
		}
		return "", &AIError{
			Message:    msg,
			StatusCode: resp.StatusCode,
		}
	}

	return string(respBody), nil
}

// GenerateTest returns a parsed AIResponse from Ollama.
func (o *OllamaClient) GenerateTest(ctx context.Context, prompt string) (AIResponse, error) {
	raw, err := o.generateRaw(ctx, prompt)
	if err != nil {
		return AIResponse{}, err
	}
	return o.parser.ParseAIResponse(raw)
}
