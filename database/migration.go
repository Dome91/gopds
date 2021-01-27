package database

import (
	"database/sql"
	"github.com/lopezator/migrator"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"path/filepath"
)

var (
	v1 = migration{"V1__init"}
)

type migration struct {
	Name  string
}

func readSQL(path string, m migration) string {
	filename := m.Name + ".sql"
	file, err := ioutil.ReadFile(filepath.Join(path, filename))
	if err != nil {
		panic(err)
	}

	return string(file)
}

func Migrate(db *sql.DB, path string) {
	m, err := migrator.New(
		migrator.Migrations(
			&migrator.Migration{
				Name: v1.Name,
				Func: func(tx *sql.Tx) error {
					if _, err := tx.Exec(readSQL(path, v1)); err != nil {
						return err
					}
					return nil
				},
			},
		),
	)

	if err != nil {
		panic(err)
	}

	err = m.Migrate(db)
	if err != nil {
		panic(err)
	}
}
