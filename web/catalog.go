package web

import (
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"gopds/services"
	"path"
	"strconv"
)

const (
	idAll         = "all"
	idDirectories = "directories"
)

type CatalogHandler struct {
	FetchAllBooksInPage                             services.FetchAllBooksInPage
	CountBooks                                      services.CountBooks
	FetchCatalogRootDirectories                     services.FetchCatalogRootDirectories
	FetchCatalogEntriesByParentCatalogEntryIDInPage services.FetchCatalogEntriesByParentCatalogEntryIDInPage
	CountCatalogEntriesByParentCatalogEntryID       services.CountCatalogEntriesByParentCatalogEntryID
	FetchCatalogEntryByID                           services.FetchCatalogEntryByID
}

func NewCatalogHandler(
	fetchAllBooksInPage services.FetchAllBooksInPage,
	countBooks services.CountBooks,
	fetchCatalogRootDirectories services.FetchCatalogRootDirectories,
	fetchCatalogEntriesByParentCatalogEntryIDInPage services.FetchCatalogEntriesByParentCatalogEntryIDInPage,
	countCatalogEntriesByParentCatalogEntryID services.CountCatalogEntriesByParentCatalogEntryID,
	fetchCatalogEntryByID services.FetchCatalogEntryByID,
) *CatalogHandler {
	return &CatalogHandler{
		FetchAllBooksInPage:         fetchAllBooksInPage,
		CountBooks:                  countBooks,
		FetchCatalogRootDirectories: fetchCatalogRootDirectories,
		FetchCatalogEntriesByParentCatalogEntryIDInPage: fetchCatalogEntriesByParentCatalogEntryIDInPage,
		CountCatalogEntriesByParentCatalogEntryID:       countCatalogEntriesByParentCatalogEntryID,
		FetchCatalogEntryByID:                           fetchCatalogEntryByID,
	}
}

func (c *CatalogHandler) Register(app *fiber.App, authorization *Authorization) {
	getPage := authorization.WithRoles(c.getPage, domain.RoleUser, domain.RoleAdmin)
	download := authorization.WithRoles(c.download, domain.RoleUser, domain.RoleAdmin)
	app.Get("/api/v1/catalog", getPage)
	app.Get("/api/v1/catalog/:id/download", download)
}

func (c *CatalogHandler) getPage(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return c.getRootPage(ctx)
	}

	page, err := strconv.Atoi(ctx.Query("page", "0"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize", "24"))
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.getPageByID(ctx, id, page, pageSize)
}

func (c *CatalogHandler) getRootPage(ctx *fiber.Ctx) error {
	allCatalogEntryResponse := catalogEntryResponse{ID: idAll, Name: "All Catalog Entries", IsDirectory: true}
	directoriesCatalogEntryResponse := catalogEntryResponse{ID: idDirectories, Name: "Catalog Directories", IsDirectory: true}
	return ctx.JSON(getCatalogEntriesResponse{CatalogEntries: []catalogEntryResponse{
		allCatalogEntryResponse,
		directoriesCatalogEntryResponse,
	}, Total: 2})
}

func (c *CatalogHandler) getPageByID(ctx *fiber.Ctx, id string, page int, pageSize int) error {
	if id == idAll {
		return c.getAllBooks(ctx, page, pageSize)
	}

	if id == idDirectories {
		return c.getRootDirectories(ctx)
	}

	entries, err := c.FetchCatalogEntriesByParentCatalogEntryIDInPage(id, page, pageSize)
	if err != nil {
		return err
	}

	count, err := c.CountCatalogEntriesByParentCatalogEntryID(id)
	if err != nil {
		return err
	}

	return ctx.JSON(getCatalogEntriesResponse{
		Total:          count,
		CatalogEntries: c.mapAllToResponse(entries),
	})
}

func (c *CatalogHandler) getRootDirectories(ctx *fiber.Ctx) error {
	directories, err := c.FetchCatalogRootDirectories()
	if err != nil {
		return err
	}

	return ctx.JSON(getCatalogEntriesResponse{
		Total:          len(directories),
		CatalogEntries: c.mapAllToResponse(directories),
	})
}

func (c *CatalogHandler) getAllBooks(ctx *fiber.Ctx, page int, pageSize int) error {
	booksInPage, err := c.FetchAllBooksInPage(page, pageSize)
	if err != nil {
		return err
	}

	count, err := c.CountBooks()
	if err != nil {
		return err
	}

	return ctx.JSON(getCatalogEntriesResponse{
		Total:          count,
		CatalogEntries: c.mapAllToResponse(booksInPage),
	})
}

func (c *CatalogHandler) download(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	catalogEntry, err := c.FetchCatalogEntryByID(id)
	if err != nil {
		return err
	}

	filename := path.Base(catalogEntry.Path)
	ctx.Set(fiber.HeaderContentDisposition, "filename="+filename)
	return ctx.SendFile(catalogEntry.Path)
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
		Type:        entry.Type,
		Cover:       entry.Cover,
	}
}

type getCatalogEntriesResponse struct {
	Total          int                    `json:"total"`
	CatalogEntries []catalogEntryResponse `json:"catalogEntries"`
}

type catalogEntryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	IsDirectory bool                    `json:"isDirectory"`
	Type        domain.CatalogEntryType `json:"type"`
	Cover       string                  `json:"cover"`
}
