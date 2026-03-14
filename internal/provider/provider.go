package provider

import (
	"context"
)

type ProviderType string

const (
	OAuthGoogle ProviderType = "google"
	OAuthGitHub ProviderType = "github"
	SAML        ProviderType = "saml"
)

type ProviderConfig struct {
	Name         string
	Type         ProviderType
	ClientID     string
	ClientSecret string
	RedirectURL  string
	// OAuth specific
	IssuerURL string
	Scopes    []string
	// SAML specific
	SAMLCert    string
	SAMLKey     string
	SAMLIDPMeta string
}

type Provider interface {
	GetName() string
	GetType() ProviderType
	AuthenticateURL(ctx context.Context, state string) (string, error)
	Callback(ctx context.Context, code string, state string) (*UserInfo, error)
}

type UserInfo struct {
	ID       string
	Email    string
	Name     string
	Provider string
}
