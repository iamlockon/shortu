package db

import (
	"context"
	"time"

	errors "github.com/iamlockon/shortu/internal/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	filter "github.com/seiflotfy/cuckoofilter"
)

type DBClient interface {
	// GetConn(context.Context) (*pg.Conn, *errors.Error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) *errors.Error

	UploadURL(c context.Context, url string, expiredAt int64, f *filter.Filter) (shorten string, err *errors.Error)
	GetURL(c context.Context, shorten string) (original string, err *errors.Error)
	LoadURL(c context.Context, f *filter.Filter) *errors.Error
}

type PgConfig struct {
	user     string
	password string
	host     string
	db       string
	timeout  int
}

// var _ DBClient = (*PgClient)(nil)

type PgxIface interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type PgClient struct {
	pool    PgxIface
	timeout time.Duration
}
