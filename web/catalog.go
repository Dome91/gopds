package web

import (
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"gopds/services"
	"strconv"
)

type CatalogHandler struct {
	fetchAllBooksInPage services.FetchAllBooksInPage
	countAllBooks       services.CountAllBooks
}

func NewCatalogHandler(
	fetchAllBooksInPage services.FetchAllBooksInPage,
	countAllBooks services.CountAllBooks,
) *CatalogHandler {
	return &CatalogHandler{
		fetchAllBooksInPage: fetchAllBooksInPage,
		countAllBooks:       countAllBooks,
	}
}

func (c *CatalogHandler) Register(app *fiber.App, authorization *Authorization) {
	getPage := authorization.WithRoles(c.getPage, domain.RoleUser, domain.RoleAdmin)
	app.Get("/api/v1/catalog", getPage)
}

func (c *CatalogHandler) getPage(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	page, err := strconv.Atoi(ctx.Query("page", "0"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "24"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	if id == "" {
		return c.getRootPage(ctx, page, pageSize)
	}

	return c.getPageByID(ctx, id, page, pageSize)
}

func (c *CatalogHandler) getRootPage(ctx *fiber.Ctx, page int, pageSize int) error {
	allCatalogEntryResponse := catalogEntryResponse{ID: "all", Name: "All Catalog Entries", IsDirectory: true}
	foldersCatalogEntryResponse := catalogEntryResponse{ID: "folders", Name: "Catalog Folders", IsDirectory: true}
	return ctx.JSON(getCatalogEntriesResponse{CatalogEntries: []catalogEntryResponse{
		allCatalogEntryResponse,
		foldersCatalogEntryResponse,
	}, Total: 2})
}

func (c *CatalogHandler) getPageByID(ctx *fiber.Ctx, id string, page int, pageSize int) error {
	if id == "all" {
		return c.getAllBooks(ctx, page, pageSize)
	}
	return nil
}

func (c *CatalogHandler) getAllBooks(ctx *fiber.Ctx, page int, pageSize int) error {
	booksInPage, err := c.fetchAllBooksInPage(page, pageSize)
	if err != nil {
		return err
	}

	count, err := c.countAllBooks()
	if err != nil {
		return err
	}

	return ctx.JSON(getCatalogEntriesResponse{
		Total:          count,
		CatalogEntries: c.mapAllToResponse(booksInPage),
	})
}

func (c *CatalogHandler) mapAllToResponse(entries []domain.CatalogEntry) []catalogEntryResponse {
	response := make([]catalogEntryResponse, len(entries))
	for index, entry := range entries {
		response[index] = c.mapToResponse(entry)
	}

	return response
}

func (c *CatalogHandler) mapToResponse(entry domain.CatalogEntry) catalogEntryResponse {
	return catalogEntryResponse{
		ID:          entry.ID,
		Name:        entry.Name,
		IsDirectory: entry.IsDirectory,
	}
}

type getCatalogEntriesResponse struct {
	Total          int                    `json:"total"`
	CatalogEntries []catalogEntryResponse `json:"catalogEntries"`
}

type catalogEntryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsDirectory bool   `json:"isDirectory"`
}
