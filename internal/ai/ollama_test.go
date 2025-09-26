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

func newOllamaClientWithMocks(parser AIParser, httpBody io.Reader, httpStatus int, httpErr error) *OllamaClient {
	return &OllamaClient{
		model:   "test-model",
		stream:  false,
		parser:  parser,
		client:  newMockHTTPClient(httpBody, httpStatus, httpErr),
		baseURL: "http://mock",
	}
}

type brokenReader struct{}

func (br *brokenReader) Read(p []byte) (int, error) { return 0, io.ErrClosedPipe }
func (br *brokenReader) Close() error               { return nil }

func TestGenerateTest_Success(t *testing.T) {
	expected := AIResponse{TestName: "T1", Code: "assert true"}
	jsonBody, _ := json.Marshal(expected)

	client := newOllamaClientWithMocks(&MockAIParser{response: expected}, strings.NewReader(string(jsonBody)), http.StatusOK, nil)
	resp, err := client.GenerateTest(context.Background(), "prompt")
	assert.NoError(t, err)
	assert.Equal(t, expected.TestName, resp.TestName)
	assert.Equal(t, expected.Code, resp.Code)
}

func TestGenerateTest_ParserError(t *testing.T) {
	parserErr := errors.New("parse error")
	client := newOllamaClientWithMocks(&MockAIParser{err: parserErr}, strings.NewReader(""), http.StatusOK, nil)
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse error")
}

func TestGenerateTest_HTTPErrorPropagation(t *testing.T) {
	client := newOllamaClientWithMocks(&MockAIParser{}, strings.NewReader(""), 0, errors.New("request failed"))
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request failed")
}

func TestGenerateTest_EmptyPrompt(t *testing.T) {
	client := newOllamaClientWithMocks(&MockAIParser{}, strings.NewReader(""), http.StatusOK, nil)
	_, err := client.GenerateTest(context.Background(), "   ") // prompt vazio
	assert.Error(t, err)
	apiErr, ok := err.(*AIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
}

func TestGenerateTest_HTTPStatusError(t *testing.T) {
	client := newOllamaClientWithMocks(&MockAIParser{}, strings.NewReader("server error"), http.StatusInternalServerError, nil)
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	apiErr, ok := err.(*AIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Contains(t, apiErr.Message, "server error")
}

func TestGenerateTest_BodyReadError(t *testing.T) {
	client := newOllamaClientWithMocks(&MockAIParser{}, &brokenReader{}, http.StatusOK, nil)
	_, err := client.GenerateTest(context.Background(), "prompt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read response body")
}
