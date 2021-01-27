package configuration

import "os"

const environmentKey = "ENV"
const environmentDevelopment = "DEV"
const migrationsKey = "MIGRATIONS"
const migrationsDefault = "migrations"
const databaseKey = "DB"
const databaseDefault = "./gopds.db"

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
	db := os.Getenv(databaseKey)
	if db == "" {
		return databaseDefault
	}

	return db
}
