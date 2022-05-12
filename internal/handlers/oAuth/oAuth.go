package oauth

import (
	"net/http"
	oauth "nwneisen/go-proxy-yourself/internal/handlers/oAuth/providers"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Provider interface used by all OAuth providers
type Provider interface {
	Begin()
	Callback()
}

// OAuth the oauth handler
type OAuth struct {
	config *config.Config
	logger *logger.Logger
}

// NewOAuth creates a new oauth handler
func NewOAuth(config *config.Config, logger *logger.Logger) *OAuth {
	return &OAuth{config, logger}
}

// ServeHTTP handles the request by passing it to the real handler
func (h *OAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Check for a valid host in the config
	host := req.Host
	_, ok := h.config.Routes[host]
	if !ok {
		h.logger.Error("Route not found in config: %s", host)
		return
	}

	// Start the authentication process
	provider := oauth.NewGoogleProvider()
	provider.Begin(w)
}
