package tokens

import (
	"io"
	"net/http"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Handlers is a generic handler for none specific routes
type Tokens struct {
	config *config.Config
	logger *logger.Logger
}

// NewCallbacks creates a new callback handler
func NewTokens(config *config.Config, logger *logger.Logger) *Tokens {
	return &Tokens{config, logger}
}

// ServeHTTP handles the request by passing it to the real handler
func (c *Tokens) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.logger.Info("Tokens handler called")
	io.WriteString(w, "Tokens received\n")

	c.logger.Info("Query values:\n")
	for key, value := range req.URL.Query() {
		c.logger.Info("%s: %s\n", key, value)
	}
}
