package database

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func withDB(f func(db *DB)) {
	db := New(":memory:")
	Migrate(db.DB.DB)
	f(db)
}

func withDBAndMock(t *testing.T, f func(db *DB, ctrl *gomock.Controller)) {
	db := New(":memory:")
	Migrate(db.DB.DB)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	f(db, ctrl)
}
