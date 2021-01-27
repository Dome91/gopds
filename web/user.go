package web

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"gopds/services"
	"net/http"
)

type UserHandler struct {
	createUser    services.CreateUser
	fetchAllUsers services.FetchAllUsers
	deleteUser    services.DeleteUser
}

func NewUserHandler(
	createUser services.CreateUser,
	fetchAllUsers services.FetchAllUsers,
	deleteUser services.DeleteUser,
) *UserHandler {
	return &UserHandler{
		createUser:    createUser,
		fetchAllUsers: fetchAllUsers,
		deleteUser:    deleteUser,
	}
}

func (u *UserHandler) Register(app *fiber.App, authorization *Authorization) {
	create := authorization.WithRoles(u.create, domain.RoleAdmin)
	getAll := authorization.WithRoles(u.getAll, domain.RoleAdmin)
	del := authorization.WithRoles(u.delete, domain.RoleAdmin)
	app.Post("/api/v1/users", create)
	app.Get("/api/v1/users", getAll)
	app.Delete("/api/v1/users/:username", del)
}

func (u *UserHandler) create(ctx *fiber.Ctx) error {
	request := createUserRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = request.Validate()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	err = u.createUser(request.Username, request.Password, request.Role)
	if err != nil {
		return err
	}

	ctx.Status(http.StatusCreated)
	return nil
}

func (u *UserHandler) getAll(ctx *fiber.Ctx) error {
	users, err := u.fetchAllUsers()
	if err != nil {
		return err
	}

	response := make([]getUserResponse, len(users))
	for index, user := range users {
		response[index] = getUserResponse{
			Username: user.Username,
			Role:     user.Role,
		}
	}

	return ctx.JSON(getAllUsersResponse{Users: response})
}

func (u *UserHandler) delete(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	if username == "" {
		return fiber.ErrBadRequest
	}

	return u.deleteUser(username)
}

type createUserRequest struct {
	Username string      `json:"username"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

func (r createUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.Role, validation.Required),
	)
}

type getAllUsersResponse struct {
	Users []getUserResponse `json:"users"`
}

type getUserResponse struct {
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
}
