package oauth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"nwneisen/go-proxy-yourself/internal/provider"
)

type GoogleProvider struct {
	config   provider.ProviderConfig
	provider *oidc.Provider
	oauthCfg oauth2.Config
}

func NewGoogleProvider(cfg provider.ProviderConfig) (*GoogleProvider, error) {
	oidcProvider, err := oidc.NewProvider(context.Background(), cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint:     oidcProvider.Endpoint(),
	}

	return &GoogleProvider{
		config:   cfg,
		provider: oidcProvider,
		oauthCfg: oauthConfig,
	}, nil
}

func (p *GoogleProvider) GetName() string {
	return p.config.Name
}

func (p *GoogleProvider) GetType() provider.ProviderType {
	return provider.OAuthGoogle
}

func (p *GoogleProvider) AuthenticateURL(ctx context.Context, state string) (string, error) {
	if state == "" {
		state = uuid.New().String()
	}
	return p.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (p *GoogleProvider) Callback(ctx context.Context, code string, state string) (*provider.UserInfo, error) {
	token, err := p.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	verifier := p.provider.Verifier(&oidc.Config{ClientID: p.config.ClientID})
	idToken, err := verifier.Verify(ctx, token.Extra("id_token").(string))
	if err != nil {
		return nil, fmt.Errorf("failed to verify id token: %w", err)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	return &provider.UserInfo{
		ID:       claims.Email,
		Email:    claims.Email,
		Name:     claims.Name,
		Provider: "google",
	}, nil
}
