package postgres

import (
	"database/sql"

	"github.com/KseniiaTD/poster/internal/database/common"
	_ "github.com/lib/pq"
)

type PostgresDatabase interface {
	common.Database
}

type postgresDB struct {
	db *sql.DB
}

func New(dsn string) (PostgresDatabase, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &postgresDB{
		db: db,
	}, nil
}

func (db *postgresDB) CloseDB() {
	db.db.Close()
}
