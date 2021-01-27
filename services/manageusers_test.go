package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"testing"
)

func TestCreateUser(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().Insert(gomock.Any()).DoAndReturn(func(user domain.User) error {
			assert.Equal(t, "username", user.Username)
			assert.Equal(t, domain.RoleUser, user.Role)
			assert.Nil(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")))
			return nil
		})
		err := CreateUserProvider(repository)("username", "password", domain.RoleUser)
		assert.Nil(t, err)
	})
}

func TestFetchAllUsers(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		user1 := domain.User{}
		user2 := domain.User{}
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().FindAll().Return([]domain.User{user1, user2}, nil)

		users, err := FetchAllUsersProvider(repository)()
		assert.Nil(t, err)
		assert.Len(t, users, 2)
		assert.Contains(t, users, user1)
		assert.Contains(t, users, user2)
	})
}

func TestFetchUserByUsername(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		user := domain.User{}
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().FindByUsername("username").Return(user, nil)

		fetchedUser, err := FetchUserByUsernameProvider(repository)("username")
		assert.Nil(t, err)
		assert.Equal(t, user, fetchedUser)
	})
}

func TestUserExistsByRole(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().ExistsByRole(domain.RoleAdmin).Return(true, nil)
		userExistsByRole, err := UserExistsByRoleProvider(repository)(domain.RoleAdmin)
		assert.Nil(t, err)
		assert.True(t, userExistsByRole)
	})
}

func TestDeleteUser(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().DeleteByUsername("user").Return(nil)
		err := DeleteUserProvider(repository)("user")
		assert.Nil(t, err)
	})
}
