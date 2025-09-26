package ai

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTest_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		prompt      string
		parser      AIParser
		httpBody    string
		httpStatus  int
		httpErr     error
		wantErr     bool
		errContains string
		wantResp    AIResponse
	}{
		{
			name:       "success",
			prompt:     "prompt",
			parser:     &MockAIParser{Response: AIResponse{TestName: "T1", Code: "assert true"}},
			httpBody:   `{"test_name":"T1","code":"assert true"}`,
			httpStatus: 200,
			wantErr:    false,
			wantResp:   AIResponse{TestName: "T1", Code: "assert true"},
		},
		{
			name:        "parser error",
			prompt:      "prompt",
			parser:      &MockAIParser{Err: errors.New("parse error")},
			httpBody:    "",
			httpStatus:  200,
			wantErr:     true,
			errContains: "parse error",
		},
		{
			name:        "http request failed",
			prompt:      "prompt",
			parser:      &MockAIParser{},
			httpBody:    "",
			httpStatus:  0,
			httpErr:     errors.New("request failed"),
			wantErr:     true,
			errContains: "request failed",
		},
		{
			name:        "empty prompt",
			prompt:      "   ",
			parser:      &MockAIParser{},
			httpBody:    "",
			httpStatus:  200,
			wantErr:     true,
			errContains: "prompt cannot be empty",
		},
		{
			name:        "http status error",
			prompt:      "prompt",
			parser:      &MockAIParser{},
			httpBody:    "server error",
			httpStatus:  500,
			wantErr:     true,
			errContains: "server error",
		},
		{
			name:        "body read error",
			prompt:      "prompt",
			parser:      &MockAIParser{},
			httpBody:    "",
			httpStatus:  200,
			httpErr:     nil,
			wantErr:     true,
			errContains: "failed to read response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var httpBodyReader io.Reader = strings.NewReader(tt.httpBody)
			if tt.name == "body read error" {
				httpBodyReader = &brokenReader{}
			}

			client := newOllamaClientWithMocks(tt.parser, httpBodyReader, tt.httpStatus, tt.httpErr)
			resp, err := client.GenerateTest(context.Background(), tt.prompt)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
