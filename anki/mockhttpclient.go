package anki

import (
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (this *MockHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := this.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (this *MockHttpClient) Get(url string) (resp *http.Response, err error) {
	args := this.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}
