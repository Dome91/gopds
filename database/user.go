package database

import (
	"gopds/domain"
	"time"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Insert(user domain.User) error {
	entity := u.mapToEntity(user)
	_, err := u.db.Exec("insert into users(username, password, role) values ($1, $2, $3)", entity.Username, entity.Password, entity.Role)
	return err
}

func (u *UserRepository) FindAll() ([]domain.User, error) {
	var entities []userEntity
	err := u.db.Select(&entities, "select * from users")
	return u.mapAllToDomain(entities), err
}

func (u *UserRepository) FindByUsername(username string) (domain.User, error) {
	var entity userEntity
	err := u.db.Get(&entity, "select * from users where username = $1", username)
	return u.mapToDomain(entity), err
}

func (u *UserRepository) ExistsByRole(role domain.Role) (bool, error) {
	var existingUsersWithRole int
	err := u.db.Get(&existingUsersWithRole, "select count(*) from users where role = $1", role)
	if err != nil {
		return false, err
	}

	return existingUsersWithRole > 0, nil
}

func (u *UserRepository) DeleteByUsername(username string) error {
	_, err := u.db.Exec("delete from users where username = $1", username)
	return err
}

func (u *UserRepository) mapToEntity(user domain.User) userEntity {
	return userEntity{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
}

func (u *UserRepository) mapAllToDomain(entities []userEntity) []domain.User {
	users := make([]domain.User, len(entities))
	for index, entity := range entities {
		users[index] = u.mapToDomain(entity)
	}

	return users
}

func (u *UserRepository) mapToDomain(entity userEntity) domain.User {
	return domain.User{
		Username: entity.Username,
		Password: entity.Password,
		Role:     entity.Role,
	}
}

type userEntity struct {
	Username  string      `db:"username"`
	Password  string      `db:"password"`
	Role      domain.Role `db:"role"`
	CreatedAt time.Time   `db:"created_at"`
}
