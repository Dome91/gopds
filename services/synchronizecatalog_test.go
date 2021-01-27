package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"path/filepath"
	"testing"
)

func TestSynchronizeCatalog(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		sourceID := "sourceID"
		path, err := filepath.Abs("../test")
		assert.Nil(t, err)

		sourceRepository := mock_domain.NewMockSourceRepository(controller)
		catalogRepository := mock_domain.NewMockCatalogRepository(controller)

		sourceRepository.EXPECT().FindByID(sourceID).Return(domain.Source{ID: sourceID, Path: path}, nil)
		catalogRepository.EXPECT().Save(gomock.Any()).DoAndReturn(func(catalog domain.Catalog) error {
			assert.Equal(t, sourceID, catalog.SourceID)

			root := catalog.Root
			assert.Equal(t, path, root.Path)
			assert.Equal(t, "test", root.Name)
			assert.True(t, root.IsDirectory)
			assert.Equal(t, domain.DIRECTORY, root.Type)
			assert.Len(t, root.Children, 1)

			ebooks := root.Children[0]
			assert.Equal(t, filepath.Join(path, "ebooks"), ebooks.Path)
			assert.Equal(t, "ebooks", ebooks.Name)
			assert.True(t, ebooks.IsDirectory)
			assert.Equal(t, domain.DIRECTORY, ebooks.Type)
			assert.Len(t, ebooks.Children, 2)

			category1 := ebooks.Children[0]
			assert.Equal(t, filepath.Join(path, "ebooks", "category1"), category1.Path)
			assert.Equal(t, "category1", category1.Name)
			assert.True(t, category1.IsDirectory)
			assert.Equal(t, domain.DIRECTORY, category1.Type)
			assert.Len(t, category1.Children, 1)

			penguinIsland := category1.Children[0]
			assert.Equal(t, filepath.Join(path, "ebooks", "category1", "Penguin Island.epub"), penguinIsland.Path)
			assert.Equal(t, "Penguin Island.epub", penguinIsland.Name)
			assert.False(t, penguinIsland.IsDirectory)
			assert.Equal(t, domain.EPUB, penguinIsland.Type)
			assert.Len(t, penguinIsland.Children, 0)

			category2 := ebooks.Children[1]
			assert.Equal(t, filepath.Join(path, "ebooks", "category2"), category2.Path)
			assert.Equal(t, "category2", category2.Name)
			assert.True(t, category2.IsDirectory)
			assert.Equal(t, domain.DIRECTORY, category2.Type)
			assert.Len(t, category2.Children, 1)

			south := category2.Children[0]
			assert.Equal(t, filepath.Join(path, "ebooks", "category2", "South.epub"), south.Path)
			assert.Equal(t, "South.epub", south.Name)
			assert.False(t, south.IsDirectory)
			assert.Equal(t, domain.EPUB, south.Type)
			assert.Len(t, south.Children, 0)
			return nil
		})

		err = SynchronizeCatalogProvider(sourceRepository, catalogRepository)(sourceID)
		assert.Nil(t, err)
	})
}
