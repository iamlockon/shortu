package db

import (
	"context"
	"fmt"
	"time"

	"github.com/iamlockon/shortu/internal/config"
	errors "github.com/iamlockon/shortu/internal/errors"
	"github.com/jackc/pgx/v4"
	pg "github.com/jackc/pgx/v4/pgxpool"
	filter "github.com/seiflotfy/cuckoofilter"
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

func (c *PgClient) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	rCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	row := c.pool.QueryRow(rCtx, sql, args...)
	return row
}

func (c *PgClient) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	rCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	rows, err := c.pool.Query(rCtx, sql, args...)
	return rows, err
}

func (c *PgClient) Exec(ctx context.Context, sql string, args ...interface{}) *errors.Error {
	rCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	commandTag, err := c.pool.Exec(rCtx, sql, args...)
	if err != nil {
		return errors.New(errors.ExecSQLError, fmt.Sprintf("failed to exec sql : %v", err))
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New(errors.ZeroAffectedSQLError, fmt.Sprintf("exec has zero affected zero, sql: %s, args: %v", sql, args))
	}
	fmt.Printf("exec sql successfully, sql :\n %s, args:\n %v\n", sql, args)
	return nil
}

func (c *PgClient) UploadURL(ctx context.Context, url string, expiredAt int64, f *filter.Filter) (string, *errors.Error) {
	sqlRecycleIfExpired := `INSERT INTO urls AS u (original, expired_at)
			VALUES (?, ?)
			ON CONFLICT (original)
			DO UPDATE SET expired_at = ?
			WHERE u.original = ? AND u.expired_at <= current_timestamp
			RETURNING id;`

	args := []interface{}{
		url,
		expiredAt,
		expiredAt,
		url,
	}
	row := c.QueryRow(ctx, sqlRecycleIfExpired, args...)
	var id int64
	if err := row.Scan(&id); err != nil {
		fmt.Println("failed to retrieve id of ", url)
		return "", errors.New(errors.QueryRowError, err.Error())
	}
	// Query id if returning id is empty
	var shorten string
	if id == 0 {
		sqlQueryIDWithOriginal := `SELECT id, shorten FROM urls WHERE original = ?`
		row := c.QueryRow(ctx, sqlQueryIDWithOriginal, url)
		err2 := row.Scan(&id, &shorten)
		if err2 != nil {
			fmt.Println("failed to retrieve id of ", url)
			return "", errors.New(errors.QueryRowError, err2.Error())
		}
	}
	if shorten == "" {
		shorten = getShortenURL(id)
		// set shorten of the row
		sqlWriteShorten := `UPDATE urls SET shorten = ? WHERE id = ?`
		if err := c.Exec(ctx, sqlWriteShorten, shorten, id); err != nil && err.Code != errors.ZeroAffectedSQLError {
			fmt.Println("failed to set shorten :", err.Msg)
			return "", err
		}
		if !f.InsertUnique([]byte(shorten)) {
			fmt.Println("shorten already exists or insert failed: ", shorten)
		}
	}
	return shorten, nil
}

func (c *PgClient) LoadURL(ctx context.Context, f *filter.Filter) *errors.Error {
	sqlLoadAllURL := `SELECT shorten FROM urls WHERE expired_at > current_timestamp;`
	rows, err := c.Query(ctx, sqlLoadAllURL)
	if err != nil {
		return errors.New(errors.QueryError, err.Error())
	}
	defer rows.Close()
	var s string
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			return errors.New(errors.ScanError, err.Error())
		}
		// set filter
		f.InsertUnique([]byte(s))
	}

	if err = rows.Err(); err != nil {
		fmt.Println("failed during iterating rows: ", err.Error())
		return errors.New(errors.RowsError, err.Error())
	}
	return nil
}

// GetURL query database for original URL, only not expired one
func (c *PgClient) GetURL(ctx context.Context, shorten string) (string, *errors.Error) {
	sqlGetURL := `SELECT original WHERE shorten = ? AND expired_at >= current_timestamp;`
	row := c.QueryRow(ctx, sqlGetURL, shorten)
	var original string
	if err := row.Scan(&original); err != nil {
		return "", errors.New(errors.URLNotFoundError, err.Error())
	}
	return original, nil
}
