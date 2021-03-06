package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"net/http"
	"testing"
)

func TestGetCatalogRootDirectories(t *testing.T) {
	handler := &CatalogHandler{}
	response, err := send(handler, "/api/v1/catalog", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var body getCatalogEntriesResponse
	parseResponse(t, response, &body)

	assert.Equal(t, 2, body.Total)
	assert.Len(t, body.CatalogEntries, 2)
	assert.Equal(t, "all", body.CatalogEntries[0].ID)
	assert.Equal(t, "All Catalog Entries", body.CatalogEntries[0].Name)
	assert.True(t, body.CatalogEntries[0].IsDirectory)
	assert.Equal(t, "directories", body.CatalogEntries[1].ID)
	assert.Equal(t, "Catalog Directories", body.CatalogEntries[1].Name)
	assert.True(t, body.CatalogEntries[1].IsDirectory)
}

func TestGetAllBooks(t *testing.T) {
	entry1 := domain.CatalogEntry{ID: "id1", Name: "name1", IsDirectory: false, Cover: "cover1"}
	entry2 := domain.CatalogEntry{ID: "id2", Name: "name2", IsDirectory: false, Cover: "cover1"}

	fetchAllBooksInPage := func(page int, pageSize int) ([]domain.CatalogEntry, error) {
		assert.Equal(t, 0, page)
		assert.Equal(t, 24, pageSize)
		return []domain.CatalogEntry{entry1, entry2}, nil
	}

	countAllBooks := func() (int, error) {
		return 20, nil
	}

	handler := &CatalogHandler{FetchAllBooksInPage: fetchAllBooksInPage, CountBooks: countAllBooks}
	response, err := send(handler, "/api/v1/catalog?id=all", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var body getCatalogEntriesResponse
	parseResponse(t, response, &body)
	assert.Equal(t, 20, body.Total)
	assert.Len(t, body.CatalogEntries, 2)
	assertCatalogEntryResponse(t, entry1, body.CatalogEntries[0])
	assertCatalogEntryResponse(t, entry2, body.CatalogEntries[1])
}

func TestGetRootDirectories(t *testing.T) {
	dir1 := domain.CatalogEntry{ID: "id1", Name: "name1", IsDirectory: true}
	dir2 := domain.CatalogEntry{ID: "id2", Name: "name2", IsDirectory: true}

	fetchCatalogRootDirectories := func() ([]domain.CatalogEntry, error) {
		return []domain.CatalogEntry{dir1, dir2}, nil
	}

	handler := &CatalogHandler{FetchCatalogRootDirectories: fetchCatalogRootDirectories}
	response, err := send(handler, "/api/v1/catalog?id=directories", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var body getCatalogEntriesResponse
	parseResponse(t, response, &body)
	assert.Equal(t, 2, body.Total)
	assert.Len(t, body.CatalogEntries, 2)
	assertCatalogEntryResponse(t, dir1, body.CatalogEntries[0])
	assertCatalogEntryResponse(t, dir2, body.CatalogEntries[1])
}

func TestGetPageByID(t *testing.T) {
	dir1 := domain.CatalogEntry{ID: "id1", Name: "name1", IsDirectory: true}
	dir2 := domain.CatalogEntry{ID: "id2", Name: "name2", IsDirectory: true}

	fetchCatalogEntriesByParentCatalogEntryIDInPage := func(parentCatalogEntryID string, page int, pageSize int) ([]domain.CatalogEntry, error) {
		assert.Equal(t, "id", parentCatalogEntryID)
		assert.Equal(t, 0, page)
		assert.Equal(t, 24, pageSize)
		return []domain.CatalogEntry{dir1, dir2}, nil
	}

	countCatalogEntriesByParentCatalogEntryID := func(parentCatalogEntryID string) (int, error) {
		return 20, nil
	}

	handler := &CatalogHandler{FetchCatalogEntriesByParentCatalogEntryIDInPage: fetchCatalogEntriesByParentCatalogEntryIDInPage, CountCatalogEntriesByParentCatalogEntryID: countCatalogEntriesByParentCatalogEntryID}
	response, err := send(handler, "/api/v1/catalog?id=id", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var body getCatalogEntriesResponse
	parseResponse(t, response, &body)
	assert.Equal(t, 20, body.Total)
	assert.Len(t, body.CatalogEntries, 2)
	assertCatalogEntryResponse(t, dir1, body.CatalogEntries[0])
	assertCatalogEntryResponse(t, dir2, body.CatalogEntries[1])
}

func TestCatalogDownload(t *testing.T) {
	fetchCatalogEntryByID := func(id string) (domain.CatalogEntry, error) {
		assert.Equal(t, "id1", id)
		return domain.CatalogEntry{Name: "South", Path: "../test/books/ebooks/epub/South.epub"}, nil
	}

	handler := &CatalogHandler{FetchCatalogEntryByID: fetchCatalogEntryByID}
	response, err := send(handler, "/api/v1/catalog/id1/download", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "filename=South.epub", response.Header.Get(fiber.HeaderContentDisposition))
	assert.Equal(t, "896191", response.Header.Get(fiber.HeaderContentLength))
}

func assertCatalogEntryResponse(t *testing.T, entry domain.CatalogEntry, response catalogEntryResponse) {
	assert.Equal(t, entry.ID, response.ID)
	assert.Equal(t, entry.Name, response.Name)
	assert.Equal(t, entry.IsDirectory, response.IsDirectory)
	assert.Equal(t, entry.Cover, response.Cover)
}
