package acquisition

import (
	"github.com/anaskhan96/soup"
)

// SoupResouceClient uses soup to get an HTML page for a ResouceClient
type SoupResouceClient struct {
}

// Get retrieves content from a webpage and returns the content or error.
func (src SoupResouceClient) Get(address string) (string, error) {
	html, err := soup.Get(address)
	if err != nil {
		return "", err
	}
	return html, nil
}
