package main

import (
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"testing"
)

func TestAdminInitializer_Initialize_createsAdmin(t *testing.T) {
	hasBeenCalled := false
	createUser := func(username string, password string, role domain.Role) error {
		assert.Equal(t, "admin", username)
		assert.Len(t, password, 8)
		assert.Equal(t, domain.RoleAdmin, role)
		hasBeenCalled = true
		return nil
	}

	userExistsByRole := func(role domain.Role) (bool, error) {
		return false, nil
	}

	NewAdminInitializer(createUser, userExistsByRole).Initialize()
	assert.True(t, hasBeenCalled)
}

func TestAdminInitializer_Initialize_doesNothingWhenAdminExists(t *testing.T) {
	createUser := func(username string, password string, role domain.Role) error {
		panic("should not be called")
	}

	userExistsByRole := func(role domain.Role) (bool, error) {
		return true, nil
	}

	NewAdminInitializer(createUser, userExistsByRole).Initialize()
}
