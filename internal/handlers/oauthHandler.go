package handlers

import (
	"nwneisen/go-proxy-yourself/internal/oauth"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

// Provider interface used by all OAuth providers
type Provider interface {
	Begin() string
	Callback()
}

// OAuth the oauth handler
type OAuthHandler struct {
	*handlers.BaseHandler
}

// NewOAuth creates a new oauth handler
func NewOAuthHandler() handlers.Handler {
	return OAuthHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// Get handles the GET request
func (oa OAuthHandler) Get() *responses.Response {

	// Check for a valid host in the config
	host := oa.Request().Host
	route, err := config.Route(host)
	if err != nil {
		return responses.NotFound("%s not found in config: %w", host, err)
	}

	logger.Info("Routing from %s to %s", host, route.EgressHostname)

	// TODO Check the query values
	logger.Info("Query values:")
	for key, value := range oa.Request().URL.Query() {
		logger.Info("%q:%q", key, value[0])
	}

	provider := oauth.NewGoogleProvider()
	return provider.Begin()

	// if req.Referer() != "https://test.nneisen.local/" {
	// 	h.googleOAuthFlow(w, req, route)
	// }

	// values := req.URL.Query()
	// if authCode, ok := values["code"]; ok {
	// 	h.googleAuthToken(w, authCode[0])
	// }
}
