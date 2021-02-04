package configuration

import (
	"os"
	"path/filepath"
)

const environmentKey = "ENV"
const environmentDevelopment = "DEV"
const migrationsKey = "MIGRATIONS"
const migrationsDefault = "migrations"
const databaseKey = "DB"
const databaseDefault = "data/gopds.db"

func IsDevelopment() bool {
	env := os.Getenv(environmentKey)
	return env == environmentDevelopment
}

func GetMigrationsPath() string {
	migrationsPath := os.Getenv(migrationsKey)
	if migrationsPath == "" {
		return migrationsDefault
	}

	return migrationsPath
}

func GetDatabasePath() string {
	createDataDirectory()
	db := os.Getenv(databaseKey)
	if db == "" {
		return databaseDefault
	}

	return db
}

func createDataDirectory() {
	dataPath := filepath.Join(".", "data")
	coversPath := filepath.Join(dataPath, "covers")

	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(coversPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
