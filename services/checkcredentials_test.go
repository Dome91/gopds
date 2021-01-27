package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"testing"
)

func TestCheckCredentials(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		hash := "$2y$10$Hzi21BlXAzynJ9Pk7iiXB.isJcp7DxTYNrnJMTdOF9FoOdkfsIptG"
		user := domain.User{Password: hash}
		repository := mock_domain.NewMockUserRepository(controller)
		repository.EXPECT().FindByUsername("username").Return(user, nil)

		err := CheckCredentialsProvider(repository)("username", "password")
		assert.Nil(t, err)
	})
}
