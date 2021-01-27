package database

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"testing"
)

func TestUserRepository_Insert(t *testing.T) {
	withDB(func(db *DB) {
		user := domain.User{Username: "username", Password: "password", Role: domain.RoleAdmin}
		err := NewUserRepository(db).Insert(user)
		assert.Nil(t, err)

		var entities []userEntity
		err = db.Select(&entities, "select * from users")
		assert.Nil(t, err)
		assert.Len(t, entities, 1)

		entity := entities[0]
		assert.Equal(t, user.Username, entity.Username)
		assert.Equal(t, user.Password, entity.Password)
		assert.Equal(t, user.Role, entity.Role)
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	withDB(func(db *DB) {
		user1 := domain.User{Username: "username1", Password: "password1", Role: domain.RoleAdmin}
		user2 := domain.User{Username: "username2", Password: "password2", Role: domain.RoleUser}

		repository := NewUserRepository(db)
		err := repository.Insert(user1)
		assert.Nil(t, err)
		err = repository.Insert(user2)
		assert.Nil(t, err)

		users, err := repository.FindAll()
		assert.Nil(t, err)
		assert.Contains(t, users, user1)
		assert.Contains(t, users, user2)
	})
}

func TestUserRepository_FindByUsername(t *testing.T) {
	withDB(func(db *DB) {
		repository := NewUserRepository(db)
		_, err := repository.FindByUsername("username")
		assert.EqualError(t, err, sql.ErrNoRows.Error())

		user := domain.User{Username: "username", Password: "password", Role: domain.RoleAdmin}
		err = repository.Insert(user)
		assert.Nil(t, err)

		foundUser, err := repository.FindByUsername("username")
		assert.Nil(t, err)
		assert.Equal(t, user.Username, foundUser.Username)
		assert.Equal(t, user.Password, foundUser.Password)
		assert.Equal(t, user.Role, foundUser.Role)
	})
}

func TestUserRepository_ExistsByRole(t *testing.T) {
	withDB(func(db *DB) {
		repository := NewUserRepository(db)
		user := domain.User{Username: "username", Password: "password", Role: domain.RoleAdmin}
		err := repository.Insert(user)
		assert.Nil(t, err)

		userExists, err := repository.ExistsByRole(domain.RoleUser)
		assert.Nil(t, err)
		assert.False(t, userExists)

		adminExists, err := repository.ExistsByRole(domain.RoleAdmin)
		assert.Nil(t, err)
		assert.True(t, adminExists)
	})
}

func TestUserRepository_DeleteByUsername(t *testing.T) {
	withDB(func(db *DB) {
		repository := NewUserRepository(db)
		user := domain.User{Username: "username", Password: "password", Role: domain.RoleAdmin}
		err := repository.Insert(user)
		assert.Nil(t, err)

		err = repository.DeleteByUsername(user.Username)
		assert.Nil(t, err)

		var count int
		err = db.Get(&count, "select count(*) from users")
		assert.Nil(t, err)
		assert.Zero(t, count)
	})
}
