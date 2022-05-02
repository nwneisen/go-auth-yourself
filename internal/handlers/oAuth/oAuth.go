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
	route, ok := h.config.Routes[host]
	if !ok {
		h.logger.Error("Route not found in config")
		return
	}

	h.logger.Info("Routing from %s to %s", host, route.EgressHostname)

	// TODO Check the query values
	h.logger.Info("Query values:")
	for key, value := range req.URL.Query() {
		h.logger.Info("%q:%q", key, value[0])
	}

	provider := oauth.NewGoogleProvider()
	provider.Begin(w)

	// if req.Referer() != "https://test.nneisen.local/" {
	// 	h.googleOAuthFlow(w, req, route)
	// }

	// values := req.URL.Query()
	// if authCode, ok := values["code"]; ok {
	// 	h.googleAuthToken(w, authCode[0])
	// }
}
