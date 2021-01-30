package web

import (
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"net/http"
	"testing"
)

func TestGetCatalogRoot(t *testing.T) {
	handler := NewCatalogHandler(nil, nil, nil)
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
	assert.Equal(t, "folders", body.CatalogEntries[1].ID)
	assert.Equal(t, "Catalog Folders", body.CatalogEntries[1].Name)
	assert.True(t, body.CatalogEntries[1].IsDirectory)
}

func TestGetAllBooks(t *testing.T) {
	entry1 := domain.CatalogEntry{ID: "id1", Name: "name1", IsDirectory: false}
	entry2 := domain.CatalogEntry{ID: "id2", Name: "name2", IsDirectory: false}

	fetchAllBooksInPage := func(page int, pageSize int) ([]domain.CatalogEntry, error) {
		assert.Equal(t, 0, page)
		assert.Equal(t, 24, pageSize)
		return []domain.CatalogEntry{entry1, entry2}, nil
	}

	countAllBooks := func() (int, error) {
		return 20, nil
	}

	handler := NewCatalogHandler(fetchAllBooksInPage, countAllBooks, nil)
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

	handler := NewCatalogHandler(nil, nil, func() ([]domain.CatalogEntry, error) {
		return []domain.CatalogEntry{dir1, dir2}, nil
	})
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

func assertCatalogEntryResponse(t *testing.T, entry domain.CatalogEntry, response catalogEntryResponse) {
	assert.Equal(t, entry.ID, response.ID)
	assert.Equal(t, entry.Name, response.Name)
	assert.Equal(t, entry.IsDirectory, response.IsDirectory)
}
