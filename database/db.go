package database

import "github.com/jmoiron/sqlx"

type DB struct {
	*sqlx.DB
}

func New(dataSourceName string) *DB {
	db := sqlx.MustOpen("sqlite3", dataSourceName)
	db.MustExec("PRAGMA foreign_keys = ON")
	return &DB{db}
}

func (db *DB) inTransaction(f func (tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
