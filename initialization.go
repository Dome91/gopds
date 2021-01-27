package main

import (
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"gopds/configuration"
	"gopds/domain"
	"gopds/services"
)

type Initializer interface {
	Initialize()
}

func Initialize(initializers ...Initializer) {
	for _, initializer := range initializers {
		initializer.Initialize()
	}
}

type AdminInitializer struct {
	createUser       services.CreateUser
	userExistsByRole services.UserExistsByRole
}

func NewAdminInitializer(createUser services.CreateUser, userExistsByRole services.UserExistsByRole) *AdminInitializer {
	return &AdminInitializer{createUser: createUser, userExistsByRole: userExistsByRole}
}

func (a *AdminInitializer) Initialize() {
	adminExists, err := a.userExistsByRole(domain.RoleAdmin)
	if err != nil {
		panic(err)
	}

	if adminExists {
		return
	}

	pass := generatePassword()
	err = a.createUser("admin", pass, domain.RoleAdmin)
	if err != nil {
		panic(err)
	}

	log.Infof("created admin with username admin and password %s", pass)
}

func generatePassword() string {
	if configuration.IsDevelopment() {
		return "pass"
	}

	pass, err := password.Generate(8, 0, 0, false, false)
	if err != nil {
		panic(err)
	}

	return pass
}
