package auth

import (
	"context"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type Provider int

const (
	GitHubAuthProvider Provider = iota
	GmailAuthProvider  Provider = iota
)

type Auth struct {
	ConfigMap  map[Provider]oauth2.Config
	signingKey string
}

func NewAuth(signingKey string, configMap map[Provider]oauth2.Config) *Auth {
	return &Auth{ConfigMap: configMap, signingKey: signingKey}
}

func (a *Auth) GetAuthCodeURL(provider Provider) string {
	config := a.ConfigMap[provider]
	return config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (a *Auth) ExchangeCode(ctx context.Context, provider Provider, code string) (*oauth2.Token, error) {
	config := a.ConfigMap[provider]
	return config.Exchange(ctx, code)
}

func (a *Auth) NewJWT(mapClaims jwt.MapClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = mapClaims
	return token.SignedString(a.signingKey)
}
