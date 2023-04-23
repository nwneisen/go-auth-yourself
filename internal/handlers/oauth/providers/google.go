package oauth

import (
	"io"
	"io/ioutil"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
	"time"
)

// GoogleProvider for a google specific provider
type GoogleProvider struct{}

// GoogleProvider constructs a new GoogleProvider structures
func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
}

// Begin starts the OAuth authentication process
func (p *GoogleProvider) Begin() *responses.Response {
	// h.logger.Info("Checking auth with Google OAuth")
	page, err := ioutil.ReadFile("web/google-redirect.html")
	if err != nil {
		// h.logger.Error(err.Error())
		return responses.InternalServerError(err.Error())
	}

	return responses.TempRedirect(string(page))
}

// Callback handles the OAuth response from Google's server
func (p *GoogleProvider) Callback(w http.ResponseWriter, authCode string) {
	// req, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", nil)

	q := req.URL.Query()
	q.Add("code", authCode)
	q.Add("client_id", "516991660211-n90f2psn5buea3n7ppucfi3iml7g1342.apps.googleusercontent.com")
	q.Add("client_secret", "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")
	q.Add("redirect_uri", "https://authed.nneisen.com")
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	req.Header = make(http.Header)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)

}
