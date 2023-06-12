package pagination

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KnightHacks/knighthacks_shared/models"
	"reflect"
)

// TODO: should cursors be encrypted? is it worth it?

var ZeroString = "0"

func DecodeCursor(cursor *string) (string, error) {
	if cursor == nil || *cursor == "" || *cursor == ZeroString {
		return ZeroString, nil
	}
	bytes, err := base64.StdEncoding.DecodeString(*cursor)
	if err != nil {
		return "", err
	}
	bytesString := string(bytes)
	return bytesString, nil
}

func GetPageInfo(a any) *models.PageInfo {
	reflectedArray := reflect.ValueOf(a)
	if reflectedArray.Kind() != reflect.Slice {
		panic("arr argument is not a slice")
	}

	if reflectedArray.Len() == 0 {
		return &models.PageInfo{
			StartCursor: ZeroString,
			EndCursor:   ZeroString,
		}
	}

	format := func(s string) string {
		bytes := []byte(s)
		return base64.StdEncoding.EncodeToString(bytes)
	}

	return &models.PageInfo{
		StartCursor: format(reflectedArray.Index(0).FieldByName("ID").String()),
		EndCursor:   format(reflectedArray.Index(reflectedArray.Len() - 1).FieldByName("ID").String()),
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
