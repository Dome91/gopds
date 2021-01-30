package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"gopds/services"
	"net/http"
	"testing"
)

func TestGetRoot(t *testing.T) {
	handler := NewOPDSHandler(services.GenerateOPDSRootFeedProvider(), nil, nil, nil, nil)
	response, err := send(handler, "/opds", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
}

func TestGetAll(t *testing.T) {
	generateOPDSAllFeed := func() (domain.Feed, error) {
		return domain.Feed{}, nil
	}
	handler := NewOPDSHandler(nil, generateOPDSAllFeed, nil, nil, nil)

	response, err := send(handler, "/opds/all", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
}

func TestGetByID(t *testing.T) {
	generateOPDSFeedByID := func(id string) (domain.Feed, error) {
		assert.Equal(t, "id1", id)
		return domain.Feed{}, nil
	}
	handler := NewOPDSHandler(nil, nil, nil, generateOPDSFeedByID, nil)

	response, err := send(handler, "/opds/id1", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
}

func TestGetFolders(t *testing.T) {
	generateOPDSFoldersFeed := func() (domain.Feed, error) {
		return domain.Feed{}, nil
	}
	handler := NewOPDSHandler(nil, nil, generateOPDSFoldersFeed, nil, nil)

	response, err := send(handler, "/opds/folders", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
}

func TestDownload(t *testing.T) {
	fetchCatalogEntryByID := func(id string) (domain.CatalogEntry, error) {
		assert.Equal(t, "id1", id)
		return domain.CatalogEntry{Name: "South.epub", Path: "../test/ebooks/epub/South.epub"}, nil
	}
	handler := NewOPDSHandler(nil, nil, nil, nil, fetchCatalogEntryByID)

	response, err := send(handler, "/opds/id1/download", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "filename=South.epub", response.Header.Get(fiber.HeaderContentDisposition))
	assert.Equal(t, "896191", response.Header.Get(fiber.HeaderContentLength))
}
