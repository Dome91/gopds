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
