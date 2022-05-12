package oauth

import (
	"io"
	"io/ioutil"
	"net/http"
)

// GoogleProvider for a google specific provider
type GoogleProvider struct{}

// GoogleProvider constructs a new GoogleProvider structures
func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}

// Begin starts the OAuth authentication process
func (p *GoogleProvider) Begin(w http.ResponseWriter) {
	page, err := ioutil.ReadFile("web/google-redirect.html")
	if err != nil {
		return
	}

	io.WriteString(w, string(page))
}
