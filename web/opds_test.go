package web

import (
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopds/database"
	"gopds/domain"
	"gopds/services"
	"gopds/util"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"testing"
)

const expectedXMLFolder = "../test/opds"

var ids = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}

func readExpectedXML(filename string) ([]byte, error) {
	filepath := path.Join(expectedXMLFolder, filename)
	file, err := ioutil.ReadFile(filepath)
	return file, err
}

func assertOPDSFeed(t *testing.T, body io.Reader, xmlFilename string) {
	actual, err := ioutil.ReadAll(body)
	assert.Nil(t, err)

	expected, err := readExpectedXML(xmlFilename)
	assert.Nil(t, err)

	var expectedFeed domain.Feed
	err = xml.Unmarshal(expected, &expectedFeed)
	assert.Nil(t, err)

	var actualFeed domain.Feed
	err = xml.Unmarshal(actual, &actualFeed)
	assert.Nil(t, err)

	assert.Equal(t, expectedFeed, actualFeed)
}

func TestGetRoot(t *testing.T) {
	handler := NewOPDSHandler(services.GenerateOPDSRootFeedProvider(), nil, nil, nil, nil)
	response, err := send(handler, "/opds", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
	assertOPDSFeed(t, response.Body, "root.xml")
}

func TestGetAll(t *testing.T) {
	withDB(func(db *database.DB) {
		setupTestCatalog(t, db)
		repository := database.NewCatalogRepository(db, util.NewSequentialIDGenerator())
		generateOPDSAllFeed := services.GenerateOPDSAllFeedProvider(repository)
		handler := NewOPDSHandler(nil, generateOPDSAllFeed, nil, nil, nil)

		response, err := send(handler, "/opds/all", http.MethodGet, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
		assertOPDSFeed(t, response.Body, "all.xml")
	})
}

func TestGetByID(t *testing.T) {
	withDB(func(db *database.DB) {
		setupTestCatalog(t, db)
		repository := database.NewCatalogRepository(db, util.NewSequentialIDGenerator())
		generateOPDSFeedByID := services.GenerateOPDSFeedByIDProvider(repository)
		handler := NewOPDSHandler(nil, nil, nil, generateOPDSFeedByID, nil)

		response, err := send(handler, "/opds/id1", http.MethodGet, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
		assertOPDSFeed(t, response.Body, "byID.xml")
	})
}

func TestGetDirectories(t *testing.T) {
	withDB(func(db *database.DB) {
		setupTestCatalog(t, db)
		repository := database.NewCatalogRepository(db, util.NewSequentialIDGenerator())
		generateOPDSDirectoriesFeed := services.GenerateOPDSDirectoriesFeedProvider(repository)
		handler := NewOPDSHandler(nil, nil, generateOPDSDirectoriesFeed, nil, nil)

		response, err := send(handler, "/opds/directories", http.MethodGet, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationXML, response.Header.Get(fiber.HeaderContentType))
		assertOPDSFeed(t, response.Body, "directories.xml")
	})
}

func TestOPDSDownload(t *testing.T) {
	fetchCatalogEntryByID := func(id string) (domain.CatalogEntry, error) {
		assert.Equal(t, "id1", id)
		return domain.CatalogEntry{Name: "South", Path: "../test/books/ebooks/epub/South.epub"}, nil
	}
	handler := NewOPDSHandler(nil, nil, nil, nil, fetchCatalogEntryByID)

	response, err := send(handler, "/opds/id1/download", http.MethodGet, nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "filename=South.epub", response.Header.Get(fiber.HeaderContentDisposition))
	assert.Equal(t, "896191", response.Header.Get(fiber.HeaderContentLength))
}

func setupTestCatalog(t *testing.T, db *database.DB) {
	checkErr := func(query string) {
		_, err := db.Exec(query)
		require.Nil(t, err)
	}

	checkErr("insert into sources(id, name, path) values('source1', 'source1', 'path/to/source')")
	checkErr("insert into catalog_entries(id, name, path, is_directory, type, parent_catalog_entry, source) values ('id1', 'dir1', 'path/dir1', true, '', null, 'source1')")
	checkErr("insert into catalog_entries(id, name, path, is_directory, type, parent_catalog_entry, source) values ('id2', 'file1.cbz', 'path/dir1/file1.cbz', false, 'CBZ', 'id1', 'source1')")
	checkErr("insert into catalog_entries(id, name, path, is_directory, type, parent_catalog_entry, source) values ('id3', 'file2.epub', 'path/dir1/file2.epub', false, 'EPUB', 'id1', 'source1')")
}
