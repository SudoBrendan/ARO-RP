package client

import "net/http"

type MockTransport struct {
	mockResponse *http.Response
	mockError    error
}

// Returns a mocked http.Client that will return the given error or response for all requests
func NewMockHttpClient(res *http.Response, err error) *http.Client {
	return &http.Client{
		Transport: NewMockTransport(res, err),
	}
}

// Returns a mocked http.Transport that will return the given error or http.Response for all requests (likely not useful directly, try NewMockHttpClient)
func NewMockTransport(res *http.Response, err error) *MockTransport {
	return &MockTransport{
		mockResponse: res,
		mockError:    err,
	}
}

// Implements http.RoundTripper interface for our mocked http.Client (likely not useful directly, try NewMockHttpClient)
func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.mockError != nil {
		return nil, t.mockError
	}
	return t.mockResponse, nil
}
