package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"io"
	"strconv"
	"time"
)

type Provider int

const (
	GitHubAuthProvider Provider = iota
	GmailAuthProvider  Provider = iota
)

type Auth struct {
	ConfigMap  map[Provider]oauth2.Config
	signingKey string
	gcm        cipher.AEAD
}

type UserClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func NewAuth(signingKey string, cipher32Bit string, configMap map[Provider]oauth2.Config) (*Auth, error) {
	newCipher, err := aes.NewCipher([]byte(cipher32Bit))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	return &Auth{ConfigMap: configMap, signingKey: signingKey, gcm: gcm}, nil
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

func (a *Auth) GetUID(ctx context.Context, provider Provider, token string) (string, error) {
	config := a.ConfigMap[provider]
	oauthClient := oauth2.NewClient(ctx, config.TokenSource(ctx, &oauth2.Token{AccessToken: token}))

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

func (a *Auth) EncryptAccessToken(token string) []byte {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	encrypted := a.gcm.Seal(nonce, nonce, []byte(token), nil)
	return encrypted
}

func (a *Auth) DecryptAccessToken(token string) ([]byte, error) {
	bytez := []byte(token)

	nonceSize := a.gcm.NonceSize()
	if len(bytez) < nonceSize {
		return nil, errors.New("size of cipher < nonce")
	}
	nonce, bytez := bytez[:nonceSize], bytez[nonceSize:]
	decryptedBytes, err := a.gcm.Open(nil, nonce, bytez, nil)
	if err != nil {
		return nil, err
	}
	return decryptedBytes, nil
}

func (a *Auth) NewTokens(userId string, role string) (refreshToken string, accessToken string, err error) {
	refreshToken, err = a.NewRefreshToken(userId, role)
	if err != nil {
		return "", "", err
	}
	accessToken, err = a.NewAccessToken(userId, role)
	if err != nil {
		return "", "", err
	}
	return refreshToken, accessToken, nil
}

func (a *Auth) NewRefreshToken(userId string, role string) (string, error) {
	return a.newJWT(userId, role, time.Minute*60)
}

func (a *Auth) NewAccessToken(userId string, role string) (string, error) {
	return a.newJWT(userId, role, time.Minute*15)
}

func (a *Auth) newJWT(userId string, role string, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := UserClaims{
		userId,
		role,
		jwt.StandardClaims{
			ExpiresAt: now.Add(expiration).UnixMilli(),
			IssuedAt:  now.UnixMilli(),
			Issuer:    "knighthacks",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(a.signingKey)
}

func (a *Auth) ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, nil
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unable to cast jwt claims to UserClaims")
	}
}
