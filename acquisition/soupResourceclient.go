package acquisition

import (
	"github.com/anaskhan96/soup"
)

type SoupResouceClient struct {
}

func (this SoupResouceClient) Get(address string) (string, error) {
	html, err := soup.Get(address)
	if err != nil {
		return "", err
	}
	return html, nil
}
