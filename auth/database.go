package auth

import (
	"context"
	"github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/models"
)

func GetUserIDByAPIKey(ctx context.Context, queryable database.Queryable, apiKey string) (userId string, role models.Role, err error) {
	err = queryable.QueryRow(
		ctx,
		"SELECT user_id, role FROM users JOIN api_keys ak on users.id = ak.user_id AND ak.key = $1",
		userId,
	).Scan(&userId, &role)
	if err != nil {
		return "", "", err
	}
	return apiKey, role, nil
}
