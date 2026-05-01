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

func (rcvr *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if rcvr.DoFunc != nil {
		return rcvr.DoFunc(req)
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

func (rcvr *MockResponseBuilder) WithStatusCode(code int) *MockResponseBuilder {
	rcvr.StatusCode = code
	return rcvr
}

func (rcvr *MockResponseBuilder) WithBody(body interface{}) *MockResponseBuilder {
	rcvr.Body = body
	return rcvr
}

func (rcvr *MockResponseBuilder) WithHeader(key, value string) *MockResponseBuilder {
	rcvr.Headers[key] = value
	return rcvr
}

func (rcvr *MockResponseBuilder) Build() *http.Response {
	var bodyReader io.ReadCloser
	if rcvr.Body != nil {
		if str, ok := rcvr.Body.(string); ok {
			bodyReader = io.NopCloser(bytes.NewBufferString(str))
		} else {
			jsonBytes, _ := json.Marshal(rcvr.Body)
			bodyReader = io.NopCloser(bytes.NewBuffer(jsonBytes))
		}
	} else {
		bodyReader = io.NopCloser(bytes.NewBufferString(""))
	}

	header := http.Header{}
	for k, v := range rcvr.Headers {
		header.Set(k, v)
	}

	return &http.Response{
		StatusCode: rcvr.StatusCode,
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
