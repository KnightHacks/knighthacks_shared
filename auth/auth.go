package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"strconv"
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
	// TODO: Implement oauth2 'state' on url to prevent CSRF https://datatracker.ietf.org/doc/html/rfc6749#section-10.12
	return config.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (a *Auth) ExchangeCode(ctx context.Context, provider Provider, code string) (*oauth2.Token, error) {
	config := a.ConfigMap[provider]
	return config.Exchange(ctx, code)
}

func (a *Auth) GetUID(ctx context.Context, provider Provider, token *oauth2.Token) (string, error) {
	config := a.ConfigMap[provider]

	oauthClient := oauth2.NewClient(ctx, config.TokenSource(ctx, token))

	if provider == GitHubAuthProvider {
		githubClient := github.NewClient(oauthClient)

		user, _, err := githubClient.Users.Get(ctx, "")
		if err != nil {
			return "", err
		}

		userId := user.ID
		if userId == nil {
			return "", errors.New("unable to retrieve github user id")
		}
		return strconv.Itoa(int(*userId)), nil
	}
	// TODO: implement for gmail
	panic("not implemented for gmail")
	return "", nil
}

func (a *Auth) NewJWT(mapClaims jwt.MapClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = mapClaims
	return token.SignedString(a.signingKey)
}
