package web

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"net/http"
	"net/http/httptest"
)

func send(handler Handler, target string, method string, body interface{}) (*http.Response, error) {
	app := fiber.New()
	withRoles := func(handler fiber.Handler, roles ...domain.Role) fiber.Handler {
		return handler
	}
	basicAuth := func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
	handler.Register(app, NewAuthorization(withRoles, basicAuth))

	m, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request := httptest.NewRequest(method, target, bytes.NewReader(m))
	request.Header.Set("Content-Type", "application/json")
	return app.Test(request)
}
