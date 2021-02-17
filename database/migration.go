package database

import (
	"database/sql"
	"embed"
	"github.com/lopezator/migrator"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

//go:embed migrations
var migrationFS embed.FS

var migrationNames = []string{"V1__init"}

func readSQL(name string) string {
	filename := name + ".sql"
	file, err := migrationFS.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		panic(err)
	}

	return string(file)
}

func Migrate(db *sql.DB) {
	var migrations []interface{}
	for _, name := range migrationNames {
		migrations = append(migrations, &migrator.Migration{
			Name: name,
			Func: func(tx *sql.Tx) error {
				_, err := tx.Exec(readSQL(name))
				return err
			}})
	}

	m, err := migrator.New(
		migrator.Migrations(migrations...),
	)

	if err != nil {
		panic(err)
	}

	err = m.Migrate(db)
	if err != nil {
		panic(err)
	}
}
