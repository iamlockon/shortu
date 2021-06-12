package db

import (
	"context"
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/config"
	errors "github.com/iamlockon/shortu/internal/errors"
	pg "github.com/jackc/pgx/v4/pgxpool"
)

// New creates a mongo client or nil if error presents
func New(cfg config.StorageConfig) (*PgClient, *errors.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetTimeout()*time.Second)
	defer cancel()

	pool, err := pg.Connect(ctx, cfg.GetConnStr())
	if err != nil {
		return nil, errors.New(errors.InvalidConfigError, fmt.Sprintf("failed to connect pg: %v", err))
	}
	DB := &PgClient{
		pool,
		cfg.GetTimeout(),
	}
	return DB, nil
}

func (c *PgClient) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

func (c *PgClient) GetConn(ctx context.Context) (*pg.Conn, *errors.Error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.New(errors.FailedToGetDBConnError, fmt.Sprintf("failed to acquire db conn : %v", err))
	}
	return conn, nil
}
