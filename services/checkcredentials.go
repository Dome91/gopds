package services

import (
	"golang.org/x/crypto/bcrypt"
	"gopds/domain"
)

type CheckCredentials func(username string, password string) error

func CheckCredentialsProvider(repository domain.UserRepository) CheckCredentials {
	return func(username string, password string) error {
		user, err := repository.FindByUsername(username)
		if err != nil {
			return err
		}

		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	}
}
