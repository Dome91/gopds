package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"testing"
)

func TestFetchCatalogEntryByID(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		entry := domain.CatalogEntry{ID: "id"}
		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindByID("id").Return(entry, nil)

		fetchedEntry, err := FetchCatalogEntryByIDProvider(repository)("id")
		assert.Nil(t, err)
		assert.Equal(t, entry, fetchedEntry)
	})
}

func TestFetchCatalogRoots(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		root1 := domain.CatalogEntry{ID: "id1"}
		root2 := domain.CatalogEntry{ID: "id2"}
		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindAllRoots().Return([]domain.CatalogEntry{root1, root2}, nil)

		roots, err := FetchCatalogRootsProvider(repository)()
		assert.Nil(t, err)
		assert.Contains(t, roots, root1)
		assert.Contains(t, roots, root2)
	})
}

func TestFetchAllBooksInPage(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		book1 := domain.CatalogEntry{ID: "id1", IsDirectory: false}
		book2 := domain.CatalogEntry{ID: "id2", IsDirectory: false}
		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindAllBooksInPage(0, 24).Return([]domain.CatalogEntry{book1, book2}, nil)

		books, err := FetchAllBooksInPageProvider(repository)(0, 24)
		assert.Nil(t, err)
		assert.Contains(t, books, book1)
		assert.Contains(t, books, book2)
	})
}

func TestCountAllBooks(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().CountAllBooks().Return(2, nil)

		count, err := CountAllBooksProvider(repository)()
		assert.Nil(t, err)
		assert.Equal(t, 2, count)
	})
}
