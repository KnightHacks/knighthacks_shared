package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
	"io"
	"strconv"
	"time"
)

type TokenType string

const (
	RefreshTokenType TokenType = "REFRESH"
	AccessTokenType  TokenType = "ACCESS"
)

var (
	TokenNotValid = errors.New("jwt token not valid")
)

type Auth struct {
	ConfigMap  map[models.Provider]oauth2.Config
	signingKey []byte
	gcm        cipher.AEAD
}

type UserClaims struct {
	UserID string      `json:"user_id"`
	Role   models.Role `json:"role"`
	Type   TokenType   `json:"type"`
	jwt.StandardClaims
}

func NewAuthWithEnvironment() (*Auth, error) {
	return NewAuth(utils.GetEnvOrDie("JWT_SIGNING_KEY"), utils.GetEnvOrDie("AES_CIPHER"), NewOAuthMap())
}

func NewAuth(signingKey string, cipher32Bit string, configMap map[models.Provider]oauth2.Config) (*Auth, error) {
	newCipher, err := aes.NewCipher([]byte(cipher32Bit))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(newCipher)
	if err != nil {
		return nil, err
	}
	return &Auth{ConfigMap: configMap, signingKey: []byte(signingKey), gcm: gcm}, nil
}

func NewOAuthMap() map[models.Provider]oauth2.Config {
	return map[models.Provider]oauth2.Config{
		//TODO: implement gmail auth, github is priority
		//auth.GmailAuthProvider: {
		//	ClientID:     "",
		//	ClientSecret: "",
		//	Endpoint:     oauth2.Endpoint{},
		//	RedirectURL:  "",
		//	Scopes:       nil,
		//},
		models.ProviderGithub: {
			ClientID:     utils.GetEnvOrDie("OAUTH_GITHUB_CLIENT_ID"),
			ClientSecret: utils.GetEnvOrDie("OAUTH_GITHUB_CLIENT_SECRET"),
			RedirectURL:  utils.GetEnvOrDie("OAUTH_GITHUB_REDIRECT_URL"),
			Endpoint:     githubOAuth.Endpoint,
			Scopes: []string{
				"read:user",
				"user:email",
			},
		},
	}
}

func (a *Auth) GetAuthCodeURL(provider models.Provider, state string) string {
	config := a.ConfigMap[provider]
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (a *Auth) ExchangeCode(ctx context.Context, provider models.Provider, code string) (*oauth2.Token, error) {
	config := a.ConfigMap[provider]
	return config.Exchange(ctx, code)
}

func (a *Auth) GetUID(ctx context.Context, provider models.Provider, token string) (string, error) {
	config := a.ConfigMap[provider]
	oauthClient := oauth2.NewClient(ctx, config.TokenSource(ctx, &oauth2.Token{AccessToken: token}))

	if provider == models.ProviderGithub {
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

func (a *Auth) NewTokens(userId string, role models.Role) (refreshToken string, accessToken string, err error) {
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

func (a *Auth) NewRefreshToken(userId string, role models.Role) (string, error) {
	return a.newJWT(userId, role, RefreshTokenType, time.Hour*24)
}

func (a *Auth) NewAccessToken(userId string, role models.Role) (string, error) {
	return a.newJWT(userId, role, AccessTokenType, time.Minute*30)
}

func (a *Auth) newJWT(userId string, role models.Role, tokenType TokenType, expiration time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := UserClaims{
		userId,
		role,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: now.Add(expiration).Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "knighthacks",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
}

func (a *Auth) ParseJWT(tokenString string, tokenType TokenType) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, TokenNotValid
	}

	if claims, ok := token.Claims.(*UserClaims); ok {
		if claims.Type != tokenType {
			return nil, fmt.Errorf("you are sending a %s token while we need a %s token", claims.Type, tokenType)
		}
		return claims, nil
	} else {
		return nil, errors.New("unable to cast jwt claims to UserClaims")
	}
}
