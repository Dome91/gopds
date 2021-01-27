package web

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUserHandler_createSucceeds(t *testing.T) {
	handler := NewUserHandler(func(username string, password string, role domain.Role) error {
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		assert.Equal(t, domain.RoleUser, role)
		return nil
	}, nil, nil)

	response, err := send(handler, "/api/v1/users", http.MethodPost, &createUserRequest{Username: "username", Password: "password", Role: domain.RoleUser})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
}

func TestUserHandler_createReturnsBadRequestForMissingFields(t *testing.T) {
	handler := NewUserHandler(func(username string, password string, role domain.Role) error {
		panic("should not be called")
	}, nil, nil)

	response, err := send(handler, "/api/v1/users", http.MethodPost, &createUserRequest{Username: "username", Role: domain.RoleUser})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	response, err = send(handler, "/api/v1/users", http.MethodPost, &createUserRequest{Password: "password", Role: domain.RoleUser})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	response, err = send(handler, "/api/v1/users", http.MethodPost, &createUserRequest{Username: "username", Password: "password"})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUserHandler_getAllSucceeds(t *testing.T) {
	user1 := domain.User{Username: "user1", Password: "pass1", Role: domain.RoleAdmin}
	user2 := domain.User{Username: "user2", Password: "pass2", Role: domain.RoleUser}

	handler := NewUserHandler(nil, func() ([]domain.User, error) {
		return []domain.User{user1, user2}, nil
	}, nil)

	response, err := send(handler, "/api/v1/users", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, response.StatusCode)

	bytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)

	var body getAllUsersResponse
	err = json.Unmarshal(bytes, &body)
	assert.Nil(t, err)

	assertUsers := func(response getUserResponse, user domain.User) {
		assert.Equal(t, user.Username, response.Username)
		assert.Equal(t, user.Role, response.Role)
	}

	assertUsers(body.Users[0], user1)
	assertUsers(body.Users[1], user2)
}


func TestUserHandler_deleteSucceeds(t *testing.T) {
	handler := NewUserHandler(nil, nil, func(username string) error {
		assert.Equal(t, "user", username)
		return nil
	})

	response, err := send(handler, "/api/v1/users/user", http.MethodDelete, nil)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, response.StatusCode)
}
