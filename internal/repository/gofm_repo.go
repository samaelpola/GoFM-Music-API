package repository

import (
	"database/sql"
	_ "embed"
	"github.com/qustavo/dotsql"
)

type DB struct {
	db  *sql.DB
	dot *dotsql.DotSql
}

//go:embed queries.sql
var queries string

func Initialize(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	dot, err := dotsql.LoadFromString(queries)
	if err != nil {
		return nil, err
	}

	_, err = dot.Exec(db, "create-database")
	if err != nil {
		return nil, err
	}

	_, err = dot.Exec(db, "create-musics-table")
	if err != nil {
		return nil, err
	}

	return &DB{
		db:  db,
		dot: dot,
	}, nil
}

func (gofmDb DB) GetSqlDb() *sql.DB {
	return gofmDb.db
}

func (gofmDb DB) GetDotSql() *dotsql.DotSql {
	return gofmDb.dot
}
