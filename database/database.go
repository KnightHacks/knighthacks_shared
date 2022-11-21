package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

// Queryable this is a sort of abstraction between pgx.Pool and Transactions, so we are able to pass either one
type Queryable interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}

func ConnectWithRetries(databaseUri string) (pool *pgxpool.Pool, err error) {
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.Connect(context.Background(), databaseUri)
		if err == nil {
			return pool, nil
		}
		time.Sleep(time.Second * 1)
	}
	return nil, err
}

func GeneratePlaceholderNumbers(start int, end int) string {
	numbers := ""
	for i := start; i <= end; i++ {
		if i == end {
			numbers += fmt.Sprintf("$%d", i)
		} else {
			numbers += fmt.Sprintf("$%d, ", i)
		}
	}
	return numbers
}
