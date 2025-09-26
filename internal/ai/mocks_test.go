package ai

import (
	"io"
	"net/http"
)

type MockAIParser struct {
	Response AIResponse
	Err      error
}

func (m *MockAIParser) ParseAIResponse(raw string) (AIResponse, error) {
	return m.Response, m.Err
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
