package pagination

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KnightHacks/knighthacks_shared/models"
)

// TODO: should cursors be encrypted? is it worth it?

var ZeroString = "0"

func DecodeCursor(cursor *string) (string, error) {
	if cursor == nil {
		return ZeroString, nil
	}
	bytes, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return "", err
	}
	bytesString := string(bytes)
	return bytesString, nil
}

func GetPageInfo(first string, last string) *models.PageInfo {
	format := func(s string) string {
		bytes := []byte(s)
		return base64.StdEncoding.EncodeToString(bytes)
	}

	return &models.PageInfo{
		StartCursor: format(first),
		EndCursor:   format(last),
	}
}

func Pagination(ctx context.Context, _ interface{}, next graphql.Resolver, maxLength int) (res interface{}, err error) {
	fieldContext := graphql.GetFieldContext(ctx)
	first := fieldContext.Args["first"].(int)
	if first > maxLength {
		return nil, fmt.Errorf("you are only allowed to request non-negative amounts less than or equal to 10")
	}
	return next(ctx)
}
