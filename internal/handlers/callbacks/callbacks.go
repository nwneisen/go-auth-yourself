package callbacks

import (
	"fmt"
	"io"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"time"
)

// Handlers is a generic handler for none specific routes
type Callbacks struct {
	config *config.Config
	logger *logger.Logger
}

// NewCallbacks creates a new callback handler
func NewCallbacks(config *config.Config, logger *logger.Logger) *Callbacks {
	return &Callbacks{config, logger}
}

// ServeHTTP handles the request by passing it to the real handler
func (c *Callbacks) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("Callback handler called")
	io.WriteString(w, "Callback received\n")

	// Get the auth code from the query values
	authCode, ok := req.URL.Query()["code"]
	if !ok {
		c.logger.Error("auth code not found in response")
	}

	// Lookup the host in the config
	host := req.Host
	route, ok := c.config.Routes[host]
	if !ok {
		c.logger.Error("Route not found in config: %s", host)
		return
	}

	// Start creating the token request
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", nil)
	if err != nil {
		c.logger.Error("error creating new request: %v", err)
	}

	// Setup the token request query values
	q := req.URL.Query()
	q.Add("code", authCode[0])
	q.Add("client_id", route.GoogleClientId)
	q.Add("client_secret", route.GoogleClientSecret)
	q.Add("redirect_uri", host+"/tokens")
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	// Setup the token request header
	req.Header = make(http.Header)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println(req)

	// Make the token request
	client := http.Client{Timeout: time.Duration(60) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
