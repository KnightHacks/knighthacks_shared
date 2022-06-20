package pagination

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
)

func Pagination(ctx context.Context, obj interface{}, next graphql.Resolver, maxLength int) (res interface{}, err error) {
	fieldContext := graphql.GetFieldContext(ctx)
	first := fieldContext.Args["first"].(int)
	if first > maxLength {
		return nil, fmt.Errorf("you are only allowed to request non-negative amounts less than or equal to 10")
	}
	return next(ctx)
}
