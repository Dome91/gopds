package web

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"gopds/services"
	"net/http"
)

type SourceHandler struct {
	createSource       services.CreateSource
	fetchAllSources    services.FetchAllSources
	deleteSource       services.DeleteSource
	synchronizeCatalog services.SynchronizeCatalog
}

func NewSourceHandler(
	createSource services.CreateSource,
	fetchAllSources services.FetchAllSources,
	deleteSource services.DeleteSource,
	synchronizeCatalog services.SynchronizeCatalog,
) *SourceHandler {
	return &SourceHandler{
		createSource:       createSource,
		fetchAllSources:    fetchAllSources,
		deleteSource:       deleteSource,
		synchronizeCatalog: synchronizeCatalog,
	}
}

func (s *SourceHandler) Register(app *fiber.App, authorization *Authorization) {
	create := authorization.WithRoles(s.create, domain.RoleAdmin)
	getAll := authorization.WithRoles(s.getAll, domain.RoleAdmin)
	del := authorization.WithRoles(s.delete, domain.RoleAdmin)
	sync := authorization.WithRoles(s.sync, domain.RoleAdmin)
	app.Post("/api/v1/sources", create)
	app.Get("/api/v1/sources", getAll)
	app.Delete("/api/v1/sources/:id", del)
	app.Put("/api/v1/sources/:id/sync", sync)
}

func (s *SourceHandler) create(ctx *fiber.Ctx) error {
	var request createSourceRequest
	err := ctx.BodyParser(&request)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = request.Validate()
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	err = s.createSource(request.Name, request.Path)
	if err != nil {
		return err
	}

	ctx.Status(http.StatusCreated)
	return nil
}

func (s *SourceHandler) getAll(ctx *fiber.Ctx) error {
	sources, err := s.fetchAllSources()
	if err != nil {
		return err
	}

	var response getAllSourcesResponse
	response.Sources = make([]getSourceResponse, len(sources))
	for index, source := range sources {
		response.Sources[index] = getSourceResponse{ID: source.ID, Name: source.Name, Path: source.Path}
	}

	return ctx.JSON(response)
}

func (s *SourceHandler) delete(ctx *fiber.Ctx) error {
	sourceID := ctx.Params("id")
	if sourceID == "" {
		return fiber.ErrBadRequest
	}

	return s.deleteSource(sourceID)
}

func (s *SourceHandler) sync(ctx *fiber.Ctx) error {
	sourceID := ctx.Params("id")
	if sourceID == "" {
		return fiber.ErrBadRequest
	}

	return s.synchronizeCatalog(sourceID)
}

type createSourceRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (r createSourceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Path, validation.Required),
	)
}

type getAllSourcesResponse struct {
	Sources []getSourceResponse `json:"sources"`
}

type getSourceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}
