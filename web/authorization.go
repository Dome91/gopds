package web

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gopds/domain"
	"gopds/services"
	"net/http"
)

const (
	roleKey = "role"
)

type Authorization struct {
	WithRoles WithRoles
	BasicAuth fiber.Handler
}

func NewAuthorization(withRoles WithRoles, basicAuth fiber.Handler) *Authorization {
	return &Authorization{WithRoles: withRoles, BasicAuth: basicAuth}
}

type WithRoles func(handler fiber.Handler, roles ...domain.Role) fiber.Handler

func WithRolesProvider(store *session.Store) WithRoles {
	return func(handler fiber.Handler, roles ...domain.Role) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			sess, err := store.Get(ctx)
			if err != nil {
				return fiber.ErrUnauthorized
			}

			role, ok := sess.Get(roleKey).(domain.Role)
			if !ok {
				return fiber.ErrUnauthorized
			}

			for _, allowedRole := range roles {
				if role == allowedRole {
					return handler(ctx)
				}
			}

			return fiber.ErrForbidden
		}
	}
}

func BasicAuthProvider(checkCredentials services.CheckCredentials) fiber.Handler {
	return basicauth.New(basicauth.Config{Authorizer: func(username string, password string) bool {
		err := checkCredentials(username, password)
		return err == nil
	}})
}

type LoginHandler struct {
	store               *session.Store
	checkCredentials    services.CheckCredentials
	fetchUserByUsername services.FetchUserByUsername
}

func NewLoginHandler(store *session.Store, checkCredentials services.CheckCredentials, fetchUserByUsername services.FetchUserByUsername) *LoginHandler {
	return &LoginHandler{store: store, checkCredentials: checkCredentials, fetchUserByUsername: fetchUserByUsername}
}

func (l *LoginHandler) Register(app *fiber.App, authorization *Authorization) {
	logout := authorization.WithRoles(l.logout, domain.RoleAdmin, domain.RoleUser)
	app.Post("/api/v1/login", l.login)
	app.Put("/api/v1/logout", logout)
}

func (l *LoginHandler) login(ctx *fiber.Ctx) error {
	var request loginRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = request.Validate()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	err = l.checkCredentials(request.Username, request.Password)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	sess, err := l.store.Get(ctx)
	if err != nil {
		return err
	}

	user, err := l.fetchUserByUsername(request.Username)
	if err != nil {
		return err
	}

	sess.Set(roleKey, user.Role)

	return sess.Save()
}

func (l *LoginHandler) logout(ctx *fiber.Ctx) error {
	sess, err := l.store.Get(ctx)
	if err != nil {
		return err
	}

	return sess.Destroy()
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r loginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
