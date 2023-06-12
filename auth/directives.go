package auth

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/golang-jwt/jwt"
	"log"
	"strings"
)

func DefaultGetUserId(ctx context.Context, obj interface{}) (string, error) {
	return "", nil
}

type HasRoleDirective struct {
	GetUserId func(ctx context.Context, obj interface{}) (string, error)
	Queryable database.Queryable
}

func (receiver HasRoleDirective) Direct(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var userClaims *UserClaims

	value, ok := ctx.Value("AuthorizationUserClaims").(*UserClaims)
	if ok {
		userClaims = value
	} else {
		auth, err := AuthFromContext(ctx)
		if err != nil {
			return nil, err
		}

		authHeader := ginContext.GetHeader("authorization")
		if len(authHeader) == 0 {
			return nil, errors.New("you must provide an authorization header")
		}

		authHeader = authHeader[7:] // remove 'bearer '

		// JWT will contain periods, API keys are alphanumeric
		if strings.Contains(authHeader, ".") {
			// 7 because it's the length of 'bearer '
			userClaims, err = auth.ParseJWT(authHeader, AccessTokenType)
			if err != nil {
				return nil, err
			}
		} else {
			userId, roleRetrieved, err := GetUserIDByAPIKey(ctx, receiver.Queryable, authHeader)
			if err != nil {
				return nil, err
			}

			userClaims = &UserClaims{
				UserID:         userId,
				Role:           roleRetrieved,
				Type:           APIKeyTokenType,
				StandardClaims: jwt.StandardClaims{},
			}

		}
	}

	if userClaims.Role == models.RoleOwns {
		return nil, errors.New("don't try to be sneaky")
	}

	switch role {
	case models.RoleAdmin:
		if userClaims.Role != models.RoleAdmin {
			return nil, errors.New("you must be an admin to use this resolver")
		}
		break
	case models.RoleNormal:
		break
	case models.RoleOwns:
		if userClaims.Role == models.RoleAdmin {
			break
		}
		id, err := receiver.GetUserId(ctx, obj)
		if err != nil {
			return nil, err
		}
		if len(id) == 0 {
			return nil, errors.New("unexpectedly the retrieved id is of length 0, possibly using DefaultGetUserId without implementing it")
		}
		log.Printf("Checking id:%s against userClaims=%v\n", id, *userClaims)
		if id != userClaims.UserID {
			return nil, errors.New("you must be own this data to use this resolver")
		}
		break
	}

	return next(context.WithValue(ctx, "AuthorizationUserClaims", userClaims))
}
