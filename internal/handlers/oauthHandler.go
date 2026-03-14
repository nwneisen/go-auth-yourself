package handlers

import (
	"context"
	"net/http"

	"nwneisen/go-proxy-yourself/internal/oauth"
	"nwneisen/go-proxy-yourself/internal/provider"
	"nwneisen/go-proxy-yourself/pkg/config"
	"nwneisen/go-proxy-yourself/pkg/logger"
	"nwneisen/go-proxy-yourself/pkg/responses"
	"nwneisen/go-proxy-yourself/pkg/server/handlers"
)

var _ = config.Route

// OAuth the oauth handler
type OAuthHandler struct {
	*handlers.BaseHandler
}

// NewOAuth creates a new oauth handler
func NewOAuth() handlers.Handler {
	return OAuthHandler{
		BaseHandler: handlers.NewBaseHandler(),
	}
}

// Get handles the GET request
func (oa OAuthHandler) Get() *responses.Response {
	host := oa.Request().Host
	route, err := config.Route(host)
	if err != nil {
		return responses.NotFound("%s not found in config: %v", host, err)
	}

	logger.Info("Routing from %s to %s", host, route.EgressHostname)

	var authProvider interface {
		AuthenticateURL(ctx context.Context, state string) (string, error)
	}

	if route.OAuth != nil {
		providerCfg := provider.ProviderConfig{
			Name:         "google",
			Type:         provider.OAuthGoogle,
			ClientID:     route.OAuth.ClientId,
			ClientSecret: route.OAuth.ClientSecret,
			RedirectURL:  route.EgressHostname + ":80/callback",
			IssuerURL:    "https://accounts.google.com",
			Scopes:       []string{"openid", "profile", "email"},
		}
		p, err := oauth.NewGoogleProvider(providerCfg)
		if err != nil {
			return responses.InternalServerError("Failed to create provider: %w", err)
		}
		authProvider = p
	} else {
		return responses.NotFound("No OAuth config found for %s", host)
	}

	authURL, err := authProvider.AuthenticateURL(context.Background(), "")
	if err != nil {
		return responses.InternalServerError("Failed to get auth URL: %w", err)
	}

	return responses.TempRedirect(authURL)
}

// Post handles the POST request
func (oa OAuthHandler) Post() *responses.Response {
	return oa.Get()
}

// Callback handles OAuth callback
func (oa OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	host := r.Host
	route, err := config.Route(host)
	if err != nil {
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	var userInfo interface{}
	if route.OAuth != nil {
		providerCfg := provider.ProviderConfig{
			Name:         "google",
			Type:         provider.OAuthGoogle,
			ClientID:     route.OAuth.ClientId,
			ClientSecret: route.OAuth.ClientSecret,
			RedirectURL:  route.EgressHostname + ":80/callback",
			IssuerURL:    "https://accounts.google.com",
			Scopes:       []string{"openid", "profile", "email"},
		}
		p, err := oauth.NewGoogleProvider(providerCfg)
		if err != nil {
			http.Error(w, "Failed to create provider", http.StatusInternalServerError)
			return
		}
		userInfo, err = p.Callback(context.Background(), code, state)
		if err != nil {
			logger.Error("OAuth callback failed: %v", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}
	}

	logger.Info("OAuth successful: %v", userInfo)
	http.Redirect(w, r, route.EgressHostname, http.StatusTemporaryRedirect)
}
