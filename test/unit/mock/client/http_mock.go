package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// MockHTTPClient implements a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"code":"SUCCESS","message":"OK"}`)),
	}, nil
}

// MockResponseBuilder helps build mock HTTP responses
type MockResponseBuilder struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

func NewMockResponseBuilder() *MockResponseBuilder {
	return &MockResponseBuilder{
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
	}
}

func (b *MockResponseBuilder) WithStatusCode(code int) *MockResponseBuilder {
	b.StatusCode = code
	return b
}

func (b *MockResponseBuilder) WithBody(body interface{}) *MockResponseBuilder {
	b.Body = body
	return b
}

func (b *MockResponseBuilder) WithHeader(key, value string) *MockResponseBuilder {
	b.Headers[key] = value
	return b
}

func (b *MockResponseBuilder) Build() *http.Response {
	var bodyReader io.ReadCloser
	if b.Body != nil {
		if str, ok := b.Body.(string); ok {
			bodyReader = io.NopCloser(bytes.NewBufferString(str))
		} else {
			jsonBytes, _ := json.Marshal(b.Body)
			bodyReader = io.NopCloser(bytes.NewBuffer(jsonBytes))
		}
	} else {
		bodyReader = io.NopCloser(bytes.NewBufferString(""))
	}

	header := http.Header{}
	for k, v := range b.Headers {
		header.Set(k, v)
	}

	return &http.Response{
		StatusCode: b.StatusCode,
		Body:       bodyReader,
		Header:     header,
	}
}

// CreateMockResponse is a convenience function to create mock responses
func CreateMockResponse(statusCode int, body interface{}) *http.Response {
	return NewMockResponseBuilder().
		WithStatusCode(statusCode).
		WithBody(body).
		Build()
}

// CreateMockJSONResponse creates a mock JSON response
func CreateMockJSONResponse(statusCode int, data interface{}) *http.Response {
	builder := NewMockResponseBuilder().
		WithStatusCode(statusCode).
		WithBody(data).
		WithHeader("Content-Type", "application/json")
	return builder.Build()
}

// CreateMockErrorResponse creates a mock error response
func CreateMockErrorResponse(statusCode int, code, message string) *http.Response {
	body := map[string]string{
		"code":    code,
		"message": message,
	}
	return CreateMockJSONResponse(statusCode, body)
}
