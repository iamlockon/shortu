package db

import (
	"context"
	"time"

	errors "github.com/iamlockon/shortu/internal/errors"
	pg "github.com/jackc/pgx/v4/pgxpool"
)

type DBClient interface {
	GetConn(context.Context) (*pg.Conn, *errors.Error)
}

type PgConfig struct {
	user     string
	password string
	host     string
	db       string
	timeout  int
}

var _ DBClient = (*PgClient)(nil)

type PgClient struct {
	pool    *pg.Pool
	timeout time.Duration
}
