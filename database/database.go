package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// Queryable this is a sort of abstraction between pgx.Pool and Transactions, so we are able to pass either one
type Queryable interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func ConnectWithRetries(databaseUri string) (pool *pgxpool.Pool, err error) {
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(context.Background(), databaseUri)
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}

	for i := 0; i < 10; i++ {
		err = pool.Ping(context.Background())
		if err != nil {
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
	return pool, err
}

func GenerateUpdatePairs(keys []any, startIndex int) string {
	pairs := ""
	for i, key := range keys {
		if i == len(keys)-1 {
			pairs += fmt.Sprintf("%s=$%d", key, i+startIndex)
		} else {
			pairs += fmt.Sprintf("%s=$%d, ", key, i+startIndex)
		}
	}
	return pairs
}
