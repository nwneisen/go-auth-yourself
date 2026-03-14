package saml

import (
	"context"
	"fmt"
	"net/url"

	"github.com/crewjam/saml"

	"nwneisen/go-proxy-yourself/internal/provider"
)

type SamlProvider struct {
	config   provider.ProviderConfig
	samlSp   *saml.ServiceProvider
	metadata []byte
	authURL  string
}

func NewSamlProvider(cfg provider.ProviderConfig) (*SamlProvider, error) {
	if cfg.SAMLIDPMeta == "" {
		return nil, fmt.Errorf("SAML IDP metadata is required")
	}

	doc := []byte(cfg.SAMLIDPMeta)
	metadataURL, _ := url.Parse(cfg.RedirectURL + "/saml/metadata")
	provider := &saml.ServiceProvider{
		MetadataURL: *metadataURL,
		EntityID:    cfg.RedirectURL + "/saml/metadata",
	}

	return &SamlProvider{
		config:   cfg,
		samlSp:   provider,
		metadata: doc,
		authURL:  cfg.RedirectURL + "/saml/SSO",
	}, nil
}

func (p *SamlProvider) GetName() string {
	return p.config.Name
}

func (p *SamlProvider) GetType() provider.ProviderType {
	return provider.SAML
}

func (p *SamlProvider) AuthenticateURL(ctx context.Context, state string) (string, error) {
	return p.authURL, nil
}

func (p *SamlProvider) Callback(ctx context.Context, code string, state string) (*provider.UserInfo, error) {
	return &provider.UserInfo{
		ID:       "saml-user-id",
		Email:    "user@example.com",
		Name:     "SAML User",
		Provider: "saml",
	}, nil
}
