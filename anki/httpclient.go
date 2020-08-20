package anki

import (
	"io"
	"net/http"
)

// HTTPClient provides a simple interface to a thing that should provide common Http functions
type HTTPClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
}
