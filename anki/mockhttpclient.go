package anki

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient mocks an HTTPClient
type MockHTTPClient struct {
	mock.Mock
}

// Post posts
func (client *MockHTTPClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := client.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Get gets
func (client *MockHTTPClient) Get(url string) (resp *http.Response, err error) {
	args := client.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}
