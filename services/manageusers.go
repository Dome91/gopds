package services

import "gopds/domain"
import "golang.org/x/crypto/bcrypt"

type CreateUser func(username string, password string, role domain.Role) error
type FetchUserByUsername func(username string) (domain.User, error)
type FetchAllUsers func() ([]domain.User, error)
type UserExistsByRole func(role domain.Role) (bool, error)
type DeleteUser func(username string) error

func CreateUserProvider(repository domain.UserRepository) CreateUser {
	return func(username string, password string, role domain.Role) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := domain.User{Username: username, Password: string(hash), Role: role}
		return repository.Insert(user)
	}
}

func FetchAllUsersProvider(repository domain.UserRepository) FetchAllUsers {
	return func() ([]domain.User, error) {
		return repository.FindAll()
	}
}

func FetchUserByUsernameProvider(repository domain.UserRepository) FetchUserByUsername {
	return func(username string) (domain.User, error) {
		return repository.FindByUsername(username)
	}
}

func UserExistsByRoleProvider(repository domain.UserRepository) UserExistsByRole {
	return func(role domain.Role) (bool, error) {
		return repository.ExistsByRole(role)
	}
}

func DeleteUserProvider(repository domain.UserRepository) DeleteUser {
	return func(username string) error {
		return repository.DeleteByUsername(username)
	}
}
