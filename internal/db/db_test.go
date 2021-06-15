package db

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func NewTestPgClient(t *testing.T) *PgClient {
	return &PgClient{
		pool:    nil,
		timeout: 1 * time.Second,
	}
}

func TestQueryRow_OK(t *testing.T) {
	pg := NewTestPgClient(t)
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	testSQL := "select id, title from books"
	expID, expTitle := 1, "one"
	rows := pgxmock.NewRows([]string{"id", "title"}).AddRow(1, "one")
	pool.ExpectQuery(testSQL).WillReturnRows(rows)
	pg.pool = pool
	if err != nil {
		t.Fatal("failed to new pgxmock pool")
	}
	actual := pg.QueryRow(context.Background(), testSQL)
	var id int
	var title string
	if err := actual.Scan(&id, &title); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expID, id)
	assert.Equal(t, expTitle, title)
}

func TestQuery_OK(t *testing.T) {
	pg := NewTestPgClient(t)
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	testSQL := "select id, title from toys"
	ids, titles := []int{1, 2}, []string{"one", "two"}
	rows := pgxmock.NewRows([]string{"id", "title"}).AddRow(ids[0], titles[0]).AddRow(ids[1], titles[1])
	pool.ExpectQuery(testSQL).WillReturnRows(rows)
	pg.pool = pool
	if err != nil {
		t.Fatal("failed to new pgxmock pool")
	}
	actual, err2 := pg.Query(context.Background(), testSQL)
	assert.Nil(t, err2)
	defer actual.Close()
	var id int
	var title string
	cnt := 0
	for actual.Next() {
		if errr := actual.Scan(&id, &title); errr != nil {
			t.Fatal(errr)
		}
		assert.Equal(t, ids[cnt], id)
		assert.Equal(t, titles[cnt], title)
		cnt++
	}
	if err = actual.Err(); err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err2)
	assert.Equal(t, 2, cnt)
}

func TestExec_OK(t *testing.T) {
	pg := NewTestPgClient(t)
	pool, err := pgxmock.NewPool()
	matcher := "INSERT INTO users"
	testSQL := "INSERT INTO users(name, created_at) VALUES (?, ?)"
	name, createdAt := "Jon", time.Now()
	pool.ExpectExec(matcher).WithArgs(name, createdAt).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	pg.pool = pool
	if err != nil {
		t.Fatal("failed to new pgxmock pool")
	}
	err2 := pg.Exec(context.Background(), testSQL, name, createdAt)
	if err2 != nil {
		t.Errorf("error '%v' was not expected, while inserting a row", err2)
	}

	if err2 := pool.ExpectationsWereMet(); err2 != nil {
		t.Errorf("there were unfulfilled expectations: %s", err2)
	}
}
