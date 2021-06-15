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

var pastTimestamp = time.Now().UTC().Add(-100 * time.Hour).Round(time.Millisecond)

// New creates a mongo client or nil if error presents
func New(cfg config.StorageConfig) (*PgClient, *errors.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetTimeout())
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

func (c *PgClient) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		fmt.Println("begin failed : ", err.Error())
		return nil, err
	}
	return tx, nil
}

func (c *PgClient) UploadURL(ctx context.Context, url string, expiredAt time.Time, f *filter.Filter) (string, *errors.Error) {
	tx, err := c.Begin(ctx)
	if err != nil {
		return "", errors.New(errors.BeginError, err.Error())
	}
	//nolint:errcheck
	defer tx.Rollback(ctx)
	sqlRecycleIfExpired := `INSERT INTO urls AS u (original, expired_at)
			VALUES ($1, $2)
			ON CONFLICT (original)
			DO UPDATE SET expired_at = $2
			WHERE u.original = $1 AND u.expired_at <= current_timestamp
			RETURNING id;`

	args := []interface{}{
		url,
		expiredAt,
	}
	row := tx.QueryRow(ctx, sqlRecycleIfExpired, args...)
	var id int64
	if err := row.Scan(&id); err != nil && err != pgx.ErrNoRows {
		fmt.Println("failed to retrieve inserted id:", err.Error())
		return "", errors.New(errors.ScanError, err.Error())
	}
	fmt.Println("[DEBUG] id: ", id)
	// Query id if record already in db
	var shorten string
	if id == 0 {
		sqlQueryIDWithOriginal := `SELECT id, shorten FROM urls WHERE original = $1`
		row := tx.QueryRow(ctx, sqlQueryIDWithOriginal, url)
		err2 := row.Scan(&id, &shorten)
		if err2 != nil {
			fmt.Println("failed to retrieve id:", err2.Error())
			return "", errors.New(errors.ScanError, err2.Error())
		}
	} else {
		// new record, let's update shorten
		shorten = getShortenURL(id)
		// set shorten of the row
		sqlWriteShorten := `UPDATE urls SET shorten = $1 WHERE id = $2`
		if cmdTag, err := tx.Exec(ctx, sqlWriteShorten, shorten, id); err != nil && cmdTag.RowsAffected() != 0 {
			fmt.Println("failed to set shorten :", err.Error())
			return "", errors.New(errors.ExecSQLError, err.Error())
		}
		if !f.InsertUnique([]byte(shorten)) {
			fmt.Println("shorten already exists or insert failed: ", shorten)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return "", errors.New(errors.CommitError, err.Error())
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

	cnt := 0
	for rows.Next() {
		var b []byte
		err = rows.Scan(&b)
		if err != nil {
			return errors.New(errors.ScanError, err.Error())
		}
		// set filter
		if b != nil {
			f.InsertUnique(b)
			cnt++
		}
	}

	if err = rows.Err(); err != nil {
		fmt.Println("failed during iterating rows: ", err.Error())
		return errors.New(errors.RowsError, err.Error())
	}
	fmt.Printf("inserted %d items to filter\n", cnt)
	return nil
}

// GetURL query database for original URL, only not expired one
func (c *PgClient) GetURL(ctx context.Context, shorten string) (string, *errors.Error) {
	sqlGetURL := `SELECT original FROM urls WHERE shorten = $1 AND expired_at >= current_timestamp;`
	row := c.QueryRow(ctx, sqlGetURL, shorten)
	var original string
	if err := row.Scan(&original); err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New(errors.URLNotFoundError, err.Error())
		}
		return "", errors.New(errors.ScanError, err.Error())
	}
	return original, nil
}

// DeleteURL set expired_at to past time so that it appears "deleted" to users.
// Note that it only sets if shorten exists in DB
func (c *PgClient) DeleteURL(ctx context.Context, shorten string) *errors.Error {
	sqlDeleteIfExists := `UPDATE urls SET expired_at = $1 WHERE shorten = $2;`
	if err := c.Exec(ctx, sqlDeleteIfExists, pastTimestamp, shorten); err != nil && err.Code != errors.ZeroAffectedSQLError {
		fmt.Println("failed to set expired_at :", err.Msg)
		return err
	}
	return nil
}
