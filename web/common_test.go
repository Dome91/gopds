package web

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gopds/database"
	"gopds/domain"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

func parseResponse(t *testing.T, response *http.Response, body interface{}) {
	bodyBytes, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(bodyBytes, body)
	assert.Nil(t, err)
}

func withDB(f func(db *database.DB)) {
	db := database.New(":memory:")
	database.Migrate(db.DB.DB)
	f(db)
}
