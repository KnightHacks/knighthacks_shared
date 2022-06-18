package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthContextMiddleware(auth *Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "Auth", auth)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AuthFromContext(ctx context.Context) (*Auth, error) {
	auth := ctx.Value("Auth")
	if auth == nil {
		err := fmt.Errorf("could not retrieve auth.Auth")
		return nil, err
	}

	if gc, ok := auth.(*Auth); ok {
		return gc, nil
	}
	return nil, errors.New("auth.Auth has wrong type")
}
