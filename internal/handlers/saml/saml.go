package saml

import (
	"net/http"
	saml "nwneisen/go-proxy-yourself/internal/handlers/saml/providers"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
)

// Provider interface used by all OAuth providers
type Provider interface {
	Begin()
	Callback()
}

type Saml struct {
	config *config.Config
	logger *logger.Logger
}

// NewSaml creates a new SAML handler
func NewSaml(config *config.Config, logger *logger.Logger) *Saml {
	return &Saml{config, logger}
}

// ServeHTTP the index handler
func (h *Saml) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.logger.Info("Saml route called")

	// Check for a valid host in the config
	host := req.Host
	_, ok := h.config.Routes[host]
	if !ok {
		h.logger.Error("Route not found in config: %s", host)
		return
	}

	// Start the authentication process
	provider := saml.NewOktaProvider()
	provider.Begin(w)
}
