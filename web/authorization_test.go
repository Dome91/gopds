package web

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"net/http"
	"testing"
)

// TODO: Test Logout
func TestLogin(t *testing.T) {
	store := session.New()
	response, err := executeLogin(t, store)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	cookie := response.Header["Set-Cookie"]
	assert.NotEmpty(t, cookie)
}

func executeLogin(t *testing.T, store *session.Store) (*http.Response, error) {
	checkCredentials := func(username string, password string) error {
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		return nil
	}

	fetchUserByUsername := func(username string) (domain.User, error) {
		assert.Equal(t, "username", username)
		return domain.User{Role: domain.RoleAdmin}, nil
	}

	handler := NewLoginHandler(store, checkCredentials, fetchUserByUsername)
	request := loginRequest{Username: "username", Password: "password"}
	response, err := send(handler, "/api/v1/login", http.MethodPost, &request)
	return response, err
}
