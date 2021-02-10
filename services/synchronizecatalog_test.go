package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	mock_util "gopds/mock/util"
	"gopds/util"
	"path/filepath"
	"testing"
)

func TestSynchronizeCatalog(t *testing.T) {
	path, err := filepath.Abs("../test/books")
	assert.Nil(t, err)

	southEpub := domain.CatalogEntry{Name: "South.epub", Path: filepath.Join(path, "ebooks", "epub", "South.epub"), IsDirectory: false, Type: domain.EPUB, Children: []domain.CatalogEntry(nil)}
	lordJimEpub := domain.CatalogEntry{Name: "Lord Jim.epub", Path: filepath.Join(path, "ebooks", "epub", "Lord Jim.epub"), IsDirectory: false, Type: domain.EPUB, Children: []domain.CatalogEntry(nil)}
	epub := domain.CatalogEntry{Name: "epub", Path: filepath.Join(path, "ebooks", "epub"), IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{lordJimEpub, southEpub}}
	lordJimAzw3 := domain.CatalogEntry{Name: "Lord Jim.azw3", Path: filepath.Join(path, "ebooks", "azw3", "Lord Jim.azw3"), IsDirectory: false, Type: domain.AZW3, Children: []domain.CatalogEntry(nil)}
	azw3 := domain.CatalogEntry{Name: "azw3", Path: filepath.Join(path, "ebooks", "azw3"), IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{lordJimAzw3}}
	ebooks := domain.CatalogEntry{Name: "ebooks", Path: filepath.Join(path, "ebooks"), IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{azw3, epub}}

	comic1 := domain.CatalogEntry{Name: "comic1.cbz", Path: filepath.Join(path, "comics", "cbz", "comic1.cbz"), IsDirectory: false, Type: domain.CBZ, Children: []domain.CatalogEntry(nil)}
	cbz := domain.CatalogEntry{Name: "cbz", Path: filepath.Join(path, "comics", "cbz"), IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{comic1}}
	comics := domain.CatalogEntry{Name: "comics", Path: filepath.Join(path, "comics"), IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{cbz}}

	root := domain.CatalogEntry{Name: "books", Path: path, IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{comics, ebooks}}

	withMock(t, func(controller *gomock.Controller) {
		sourceID := "sourceID"

		sourceRepository := mock_domain.NewMockSourceRepository(controller)
		catalogRepository := mock_domain.NewMockCatalogRepository(controller)

		sourceRepository.EXPECT().FindByID(sourceID).Return(domain.Source{ID: sourceID, Path: path}, nil)
		catalogRepository.EXPECT().Save(gomock.Any()).DoAndReturn(func(catalog domain.Catalog) error {
			assert.Equal(t, sourceID, catalog.SourceID)
			assert.Equal(t, root, catalog.Root)
			return nil
		})

		catalogRepository.EXPECT().FindAllBooksWithoutCover().DoAndReturn(func() ([]domain.CatalogEntry, error) {
			savedSouthEpub := southEpub
			savedSouthEpub.ID = "id1"
			return []domain.CatalogEntry{savedSouthEpub}, nil
		})

		bus := mock_util.NewMockBus(controller)
		bus.EXPECT().Emit(gomock.Any(), util.GenerateCoverTopic, gomock.AssignableToTypeOf(util.GenerateCoverEvent{})).Times(1)

		err = SynchronizeCatalogProvider(sourceRepository, catalogRepository, bus)(sourceID)
		assert.Nil(t, err)
	})
}
