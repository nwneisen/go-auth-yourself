package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"nwneisen/go-proxy-yourself/internal/provider"
)

type GitHubProvider struct {
	config   provider.ProviderConfig
	oauthCfg oauth2.Config
}

func NewGitHubProvider(cfg provider.ProviderConfig) (*GitHubProvider, error) {
	oauthConfig := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	return &GitHubProvider{
		config:   cfg,
		oauthCfg: oauthConfig,
	}, nil
}

func (p *GitHubProvider) GetName() string {
	return p.config.Name
}

func (p *GitHubProvider) GetType() provider.ProviderType {
	return provider.OAuthGitHub
}

func (p *GitHubProvider) AuthenticateURL(ctx context.Context, state string) (string, error) {
	if state == "" {
		state = uuid.New().String()
	}
	return p.oauthCfg.AuthCodeURL(state), nil
}

func (p *GitHubProvider) Callback(ctx context.Context, code string, state string) (*provider.UserInfo, error) {
	token, err := p.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := p.oauthCfg.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userData struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, fmt.Errorf("failed to decode user data: %w", err)
	}

	return &provider.UserInfo{
		ID:       fmt.Sprintf("%d", userData.ID),
		Email:    userData.Email,
		Name:     userData.Name,
		Provider: "github",
	}, nil
}
