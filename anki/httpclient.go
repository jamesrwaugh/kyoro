package anki

import (
	"io"
	"net/http"
)

type HttpClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
}
