package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockAIParser struct {
	response AIResponse
	err      error
}

func (m *MockAIParser) ParseAIResponse(raw string) (AIResponse, error) {
	return m.response, m.err
}

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func newMockHTTPClient(body io.Reader, status int, err error) *MockHTTPClient {
	return &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(body),
			}, nil
		},
	}
}

func newOllamaClientWithMocks(resp AIResponse, parseErr error, httpBody io.Reader, httpStatus int, httpErr error) *OllamaClient {
	return &OllamaClient{
		model:   "test-model",
		stream:  false,
		parser:  &MockAIParser{response: resp, err: parseErr},
		client:  newMockHTTPClient(httpBody, httpStatus, httpErr),
		baseURL: "http://mock",
	}
}

func TestGenerateRaw_HTTPError(t *testing.T) {
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(""), 0, io.ErrUnexpectedEOF)
	_, err := client.GenerateRaw(context.Background(), "Test prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestGenerateRaw_EmptyPrompt(t *testing.T) {
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(""), 0, nil)
	_, err := client.GenerateRaw(context.Background(), "   ")
	assert.Error(t, err)
	apiErr, ok := err.(*AIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
}

func TestGenerateRaw_APIErrorStatus(t *testing.T) {
	errorMsg := "Bad Request"
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(errorMsg), http.StatusInternalServerError, nil)
	_, err := client.GenerateRaw(context.Background(), "Test prompt")
	apiErr, ok := err.(*AIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, errorMsg, apiErr.Message)
}

func TestGenerateRaw_BodyReadError(t *testing.T) {
	client := newOllamaClientWithMocks(AIResponse{}, nil, &brokenReader{}, http.StatusOK, nil)
	_, err := client.GenerateRaw(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read response body")
}

func TestGenerateRaw_LargeErrorTruncated(t *testing.T) {
	largeBody := strings.Repeat("x", 3000)
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(largeBody), http.StatusInternalServerError, nil)

	_, err := client.GenerateRaw(context.Background(), "prompt")
	apiErr, ok := err.(*AIError)
	assert.True(t, ok)
	assert.Contains(t, apiErr.Message, "[truncated]")
}

func TestGenerateRaw_Success(t *testing.T) {
	mockResp := `{"test_name":"Test1","code":"assert true"}`
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(mockResp), http.StatusOK, nil)
	raw, err := client.GenerateRaw(context.Background(), "prompt text")
	assert.NoError(t, err)
	assert.Equal(t, mockResp, raw)
}

func TestGenerateTest_Success(t *testing.T) {
	parser := &MockAIParser{
		response: AIResponse{TestName: "T1", Code: "assert true"},
	}
	jsonBody, _ := json.Marshal(parser.response)
	client := newOllamaClientWithMocks(parser.response, nil, strings.NewReader(string(jsonBody)), http.StatusOK, nil)
	resp, err := client.GenerateTest(context.Background(), "prompt")
	assert.NoError(t, err)
	assert.Equal(t, "T1", resp.TestName)
}

func TestGenerateTest_ParserError(t *testing.T) {
	parser := &MockAIParser{err: errors.New("parse error")}
	client := newOllamaClientWithMocks(AIResponse{}, parser.err, strings.NewReader(""), http.StatusOK, nil)
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse error")
}

func TestGenerateTest_PropagatesGenerateRawError(t *testing.T) {
	client := newOllamaClientWithMocks(AIResponse{}, nil, strings.NewReader(""), 0, errors.New("request failed"))
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

type brokenReader struct{}

func (br *brokenReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func (br *brokenReader) Close() error {
	return nil
}
