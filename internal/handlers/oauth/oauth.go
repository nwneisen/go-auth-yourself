package oauth

import (
	oauth "nwneisen/go-proxy-yourself/internal/handlers/oauth/providers"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
	"nwneisen/go-proxy-yourself/pkg/server/responses"
)

// Provider interface used by all OAuth providers
type Provider interface {
	Begin() string
	Callback()
}

// OAuth the oauth handler
type OAuth struct {
	*handlers.BaseHandler
}

// NewOAuth creates a new oauth handler
func NewOAuth(config *config.Config, logger *logger.Logger) handlers.Handler {
	return OAuth{
		BaseHandler: handlers.NewBaseHandler(config, logger),
	}
}

// Get handles the GET request
func (oa OAuth) Get() *responses.Response {

	// Check for a valid host in the config
	host := oa.Request().Host
	route, ok := oa.Config().Routes[host]
	if !ok {
		oa.Log().Error("Route not found in config: %s", host)
		return responses.NotFound("Route not found in config")
	}

	oa.Log().Info("Routing from %s to %s", host, route.EgressHostname)

	// TODO Check the query values
	oa.Log().Info("Query values:")
	for key, value := range oa.Request().URL.Query() {
		oa.Log().Info("%q:%q", key, value[0])
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
