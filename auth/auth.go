package auth

import (
	"context"
	"golang.org/x/oauth2"
)

type Provider int

const (
	GitHubAuthProvider Provider = iota
	GmailAuthProvider  Provider = iota
)

type Auth struct {
	configMap map[Provider]oauth2.Config
}

func (a Auth) GetAuthCodeURL(provider Provider) string {
	config := a.configMap[provider]
	return config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (a Auth) ExchangeCode(ctx context.Context, provider Provider, code string) (*oauth2.Token, error) {
	config := a.configMap[provider]
	return config.Exchange(ctx, code)
}
