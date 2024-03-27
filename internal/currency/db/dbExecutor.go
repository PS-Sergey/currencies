package db

import (
	"context"
	"database/sql"
)

type Executor interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
