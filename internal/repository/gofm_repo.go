package repository

import (
	"database/sql"
	"github.com/qustavo/dotsql"
)

type DB struct {
	db  *sql.DB
	dot *dotsql.DotSql
}

func Initialize(driverName string) (*DB, error) {
	db, err := sql.Open(driverName, "file:internal/database/gofm.sqlite")
	if err != nil {
		return nil, err
	}

	dot, err := dotsql.LoadFromFile("internal/repository/migration.sql")
	if err != nil {
		return nil, err
	}

	_, err = dot.Exec(db, "create-database")
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
