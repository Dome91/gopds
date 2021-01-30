package services

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	mock_domain "gopds/mock/domain"
	"testing"
)

func TestGenerateOPDSRootFeed(t *testing.T) {
	feed := GenerateOPDSRootFeedProvider()()
	assert.Equal(t, "GOPDS", feed.ID)
	assert.Equal(t, "Root", feed.Title)

	links := feed.Links
	assert.Len(t, links, 2)
	assert.Equal(t, "start", links[0].Rel)
	assert.Equal(t, "/opds", links[0].Href)
	assert.EqualValues(t, "application/atom+xml;profile=opds-catalog;kind=acquisition", links[0].Type)
	assert.Equal(t, "self", links[1].Rel)
	assert.Equal(t, "/opds", links[1].Href)
	assert.EqualValues(t, "application/atom+xml;profile=opds-catalog;kind=acquisition", links[1].Type)

	entries := feed.Entries
	assert.Len(t, entries, 2)

	entry1 := entries[0]
	assert.Equal(t, "All Catalog Entries", entry1.Title)
	entry1Links := entry1.Links
	assert.Len(t, entry1Links, 1)
	assert.Equal(t, "http://opds-spec.org/crawlable", entry1Links[0].Rel)
	assert.Equal(t, "/opds/all", entry1Links[0].Href)
	assert.EqualValues(t, "application/atom+xml;profile=opds-catalog;kind=acquisition", entry1Links[0].Type)

	entry2 := entries[1]
	assert.Equal(t, "Catalog Directories", entry2.Title)
	entry2Links := entry2.Links
	assert.Len(t, entry1Links, 1)
	assert.Equal(t, "http://opds-spec.org/crawlable", entry2Links[0].Rel)
	assert.Equal(t, "/opds/directories", entry2Links[0].Href)
	assert.EqualValues(t, "application/atom+xml;profile=opds-catalog;kind=acquisition", entry2Links[0].Type)
}

func TestGenerateOPDSAllFeed(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		catalogEntry1 := domain.CatalogEntry{ID: "id1", Name: "catalogEntry1", Type: domain.EPUB}
		catalogEntry2 := domain.CatalogEntry{ID: "id2", Name: "catalogEntry2", Type: domain.CBZ}

		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindAllBooks().Return([]domain.CatalogEntry{catalogEntry1, catalogEntry2}, nil)

		feed, err := GenerateOPDSAllFeedProvider(repository)()
		assert.Nil(t, err)
		assert.Equal(t, "GOPDS", feed.ID)
		assert.Equal(t, "All Catalog Entries", feed.Title)

		entries := feed.Entries
		assert.Len(t, entries, 2)
		assert.Equal(t, catalogEntry1.ID, entries[0].ID)
		assert.Equal(t, catalogEntry1.Name, entries[0].Title)
		entry1Links := entries[0].Links
		assert.Len(t, entry1Links, 1)
		assert.Equal(t, "/opds/id1/download", entry1Links[0].Href)
		assert.Equal(t, "application/epub+zip", string(entry1Links[0].Type))

		assert.Equal(t, catalogEntry2.ID, entries[1].ID)
		assert.Equal(t, catalogEntry2.Name, entries[1].Title)
		entry2Links := entries[1].Links
		assert.Len(t, entry2Links, 1)
		assert.Equal(t, "/opds/id2/download", entry2Links[0].Href)
		assert.Equal(t, "application/vnd.comicbook+zip", string(entry2Links[0].Type))
	})
}

func TestGenerateOPDSDirectoriesFeed(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		catalogEntry1 := domain.CatalogEntry{ID: "id1", Name: "catalogEntry1"}
		catalogEntry2 := domain.CatalogEntry{ID: "id2", Name: "catalogEntry2"}

		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindAllRoots().Return([]domain.CatalogEntry{catalogEntry1, catalogEntry2}, nil)

		feed, err := GenerateOPDSDirectoriesFeedProvider(repository)()
		assert.Nil(t, err)
		assert.Equal(t, "GOPDS", feed.ID)
		assert.Equal(t, "Catalog Directories", feed.Title)

		entries := feed.Entries
		assert.Len(t, entries, 2)
		assert.Equal(t, catalogEntry1.ID, entries[0].ID)
		assert.Equal(t, catalogEntry1.Name, entries[0].Title)
		entry1Links := entries[0].Links
		assert.Len(t, entry1Links, 1)
		assert.Equal(t, "/opds/id1", entry1Links[0].Href)

		assert.Equal(t, catalogEntry2.ID, entries[1].ID)
		assert.Equal(t, catalogEntry2.Name, entries[1].Title)
		entry2Links := entries[1].Links
		assert.Len(t, entry2Links, 1)
		assert.Equal(t, "/opds/id2", entry2Links[0].Href)
	})
}

func TestGenerateOPDSFeedByID(t *testing.T) {
	withMock(t, func(controller *gomock.Controller) {
		catalogEntry1 := domain.CatalogEntry{ID: "id1", Name: "catalogEntry1"}
		catalogEntry2 := domain.CatalogEntry{ID: "id2", Name: "catalogEntry2", IsDirectory: true}
		catalogEntry3 := domain.CatalogEntry{ID: "id3", Name: "catalogEntry3", IsDirectory: false}

		repository := mock_domain.NewMockCatalogRepository(controller)
		repository.EXPECT().FindByID("id1").Return(catalogEntry1, nil)
		repository.EXPECT().FindAllByParentCatalogEntryID("id1").Return([]domain.CatalogEntry{catalogEntry2, catalogEntry3}, nil)

		feed, err := GenerateOPDSFeedByIDProvider(repository)("id1")
		assert.Nil(t, err)
		assert.Equal(t, catalogEntry1.Name, feed.Title)
		assert.Len(t, feed.Entries, 2)

		assert.Equal(t, catalogEntry2.Name, feed.Entries[0].Title)
		assert.Equal(t, catalogEntry2.ID, feed.Entries[0].ID)
		assert.Len(t, feed.Entries[0].Links, 1)
		assert.Equal(t, "/opds/id2", feed.Entries[0].Links[0].Href)

		assert.Equal(t, catalogEntry3.Name, feed.Entries[1].Title)
		assert.Equal(t, catalogEntry3.ID, feed.Entries[1].ID)
		assert.Len(t, feed.Entries[1].Links, 1)
		assert.Equal(t, "/opds/id3/download", feed.Entries[1].Links[0].Href)
	})
}
