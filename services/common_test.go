package services

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func withMock(t *testing.T, f func(controller *gomock.Controller)) {
	controller := gomock.NewController(t)
	f(controller)
	controller.Finish()
}
