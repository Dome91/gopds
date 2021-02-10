package services

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/mustafaturan/bus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"gopds/util"
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCoverForCBZ(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		entry := domain.CatalogEntry{ID: "id1", Name: "comic1.cbz", Path: "../test/books/comics/cbz/comic1.cbz", Type: domain.CBZ}
		var generatedCover string

		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindByID(entry.ID).Return(entry, nil)
		repository.EXPECT().UpdateSetCoverByID(entry.ID, gomock.Any()).DoAndReturn(func(_ string, cover string) error {
			ext := filepath.Ext(cover)
			assert.Equal(t, ".jpg", ext)

			var name = cover[0 : len(cover)-len(ext)]
			_, err := uuid.Parse(name)
			assert.Nil(t, err)

			generatedCover = cover
			return nil
		})

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("data/covers", os.ModePerm)
		assert.Nil(t, err)

		GenerateCoverProvider(fs, repository)(&bus.Event{Data: util.GenerateCoverEvent{ID: entry.ID}})

		coverFile, err := fs.Open("data/covers/" + generatedCover)
		require.Nil(t, err)
		defer coverFile.Close()

		img, _, err := image.Decode(coverFile)
		assert.Nil(t, err)
		assert.Equal(t, 800, img.Bounds().Max.Y)
	})
}
