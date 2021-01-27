package database

func withDB(f func(db *DB)) {
	db := New(":memory:")
	Migrate(db.DB.DB, "../migrations")
	f(db)
}
