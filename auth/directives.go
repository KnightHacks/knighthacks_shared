package auth

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/utils"
)

type HasRoleDirective struct {
	GetId func(ctx context.Context, obj interface{}) (string, error)
}

func (receiver HasRoleDirective) Direct(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
	ginContext, err := utils.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := AuthFromContext(ctx)
	if err != nil {
		return nil, err
	}

	authHeader := ginContext.GetHeader("authorization")

	userClaims, err := auth.ParseJWT(authHeader)
	if err != nil {
		return nil, err
	}
	/*
		if userClaims.Role == auth.RoleOwns {
			return nil, errors.New("don't try to be sneaky")
		}

		switch role {
		case auth.RoleAdmin:
			if authorizationContext.Role != auth.RoleAdmin {
				return nil, errors.New("you must be an admin to use this resolver")
			}
		case auth.RoleDefault:
			break
		case auth.RoleOwns:
			if authorizationContext.Role == auth.RoleAdmin {
				break
			}
			id, err := receiver.GetId(ctx, obj)
			if err != nil {
				return nil, err
			}
			log.Printf("Checking id:%s against authorizationContext:%s\n", id, authorizationContext)
			if id != authorizationContext.ID {
				return nil, errors.New("you must be own this data to use this resolver")
			}
		case auth.RoleServer:
			if authorizationContext.Role != auth.RoleServer {
				return nil, errors.New("you cannot call this service directly")
			}
		}
	*/

	return next(context.WithValue(ctx, "AuthorizationUserClaims", userClaims))
}
