package database

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopds/domain"
	"gopds/util"
	"testing"
)

func TestCatalogRepository_Save(t *testing.T) {
	withDBAndMock(t, func(db *DB, ctrl *gomock.Controller) {
		book1, book2, book3, ebooks, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var entities []catalogEntryEntity
		err = db.Select(&entities, "select * from catalog_entries order by path")
		assert.Nil(t, err)
		assert.Len(t, entities, 3)

		ebooksEntity := entities[0]
		assertCatalogEntry(t, ebooksEntity, ebooks, catalog.SourceID)
		assert.False(t, ebooksEntity.ParentCatalogEntryID.Valid)

		bookEntity1 := entities[1]
		assertCatalogEntry(t, bookEntity1, book1, catalog.SourceID)
		assert.Equal(t, ebooksEntity.ID, bookEntity1.ParentCatalogEntryID.String)

		bookEntity2 := entities[2]
		assertCatalogEntry(t, bookEntity2, book2, catalog.SourceID)
		assert.Equal(t, ebooksEntity.ID, bookEntity2.ParentCatalogEntryID.String)

		ebooks.Children = []domain.CatalogEntry{book1, book3}
		catalog.Root = ebooks
		err = repository.Save(catalog)
		assert.Nil(t, err)

		entities = nil
		err = db.Select(&entities, "select * from catalog_entries order by path")
		assert.Nil(t, err)
		assert.Len(t, entities, 3)

		ebooksEntity = entities[0]
		assertCatalogEntry(t, ebooksEntity, ebooks, catalog.SourceID)
		assert.False(t, ebooksEntity.ParentCatalogEntryID.Valid)

		bookEntity1 = entities[1]
		assertCatalogEntry(t, bookEntity1, book1, catalog.SourceID)
		assert.Equal(t, ebooksEntity.ID, bookEntity1.ParentCatalogEntryID.String)

		bookEntity3 := entities[2]
		assertCatalogEntry(t, bookEntity3, book3, catalog.SourceID)
		assert.Equal(t, ebooksEntity.ID, bookEntity3.ParentCatalogEntryID.String)
	})
}

func assertCatalogEntry(t *testing.T, entity catalogEntryEntity, catalogEntry domain.CatalogEntry, sourceID string) {
	assert.Equal(t, catalogEntry.Name, entity.Name)
	assert.Equal(t, catalogEntry.Path, entity.Path)
	assert.NotNil(t, entity.CreatedAt)
	assert.Equal(t, catalogEntry.IsDirectory, entity.IsDirectory)
	assert.Equal(t, catalogEntry.Type, entity.Type)
	assert.True(t, entity.FoundDuringLastSync)
	assert.Equal(t, sourceID, entity.SourceID)
}

func TestCatalogRepository_FindAllBooks(t *testing.T) {
	withDB(func(db *DB) {
		book1, book2, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		entries, err := repository.FindAllBooks()
		assert.Nil(t, err)
		assert.Len(t, entries, 2)

		entry1 := entries[0]
		assert.Equal(t, book1.Name, entry1.Name)
		assert.Equal(t, book1.Path, entry1.Path)
		assert.False(t, entry1.IsDirectory)

		entry2 := entries[1]
		assert.Equal(t, book2.Name, entry2.Name)
		assert.Equal(t, book2.Path, entry2.Path)
		assert.False(t, entry2.IsDirectory)
	})
}

func TestCatalogRepository_FindByID(t *testing.T) {
	withDB(func(db *DB) {
		_, _, _, _, catalog := generateCatalog(db)
		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var entity catalogEntryEntity
		err = db.Get(&entity, "select * from catalog_entries limit 1")
		assert.Nil(t, err)

		catalogEntry, err := repository.FindByID(entity.ID)
		assert.Nil(t, err)

		assertCatalogEntry(t, entity, catalogEntry, catalog.SourceID)
	})
}

func TestCatalogRepository_FindAllRootDirectories(t *testing.T) {
	withDB(func(db *DB) {
		_, _, _, ebooks, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		roots, err := repository.FindAllRootDirectories()
		assert.Nil(t, err)
		assert.Len(t, roots, 1)

		ebooks.ID = roots[0].ID
		ebooks.Children = nil
		assert.Equal(t, ebooks, roots[0])
	})
}

func TestCatalogRepository_FindAllByParentCatalogEntryID(t *testing.T) {
	withDB(func(db *DB) {
		book1, book2, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var parentEntity catalogEntryEntity
		err = db.Get(&parentEntity, "select * from catalog_entries where parent_catalog_entry is null")
		assert.Nil(t, err)

		children, err := repository.FindAllByParentCatalogEntryID(parentEntity.ID)
		assert.Nil(t, err)
		assert.Len(t, children, 2)

		book1.ID = children[0].ID
		book2.ID = children[1].ID
		assert.Equal(t, book1, children[0])
		assert.Equal(t, book2, children[1])
	})
}

func TestCatalogRepository_FindAllByParentCatalogEntryIDInPage(t *testing.T) {
	withDB(func(db *DB) {
		book1, _, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var parentEntity catalogEntryEntity
		err = db.Get(&parentEntity, "select * from catalog_entries where parent_catalog_entry is null")
		assert.Nil(t, err)

		children, err := repository.FindAllByParentCatalogEntryIDInPage(parentEntity.ID, 0, 1)
		assert.Nil(t, err)
		assert.Len(t, children, 1)

		book1.ID = children[0].ID
		assert.Equal(t, book1, children[0])
	})
}

func TestCatalogRepository_FindAllBooksInPage(t *testing.T) {
	withDB(func(db *DB) {
		book1, _, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		booksInPage, err := repository.FindAllBooksInPage(0, 1)
		assert.Nil(t, err)

		book1.ID = booksInPage[0].ID
		assert.Equal(t, book1, booksInPage[0])
	})
}

func TestCatalogRepository_CountBooks(t *testing.T) {
	withDB(func(db *DB) {
		_, _, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		count, err := repository.CountBooks()
		assert.Nil(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestCatalogRepository_CountByParentCatalogEntryID(t *testing.T) {
	withDB(func(db *DB) {
		_, _, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var parentEntity catalogEntryEntity
		err = db.Get(&parentEntity, "select * from catalog_entries where parent_catalog_entry is null")
		assert.Nil(t, err)

		count, err := repository.CountByParentCatalogEntryID(parentEntity.ID)
		assert.Nil(t, err)
		assert.Equal(t, 2, count)
	})
}

func TestCatalogRepository_UpdateSetCoverByID(t *testing.T) {
	withDB(func(db *DB) {
		_, _, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		var entity catalogEntryEntity
		err = db.Get(&entity, "select * from catalog_entries where parent_catalog_entry is null")
		assert.Nil(t, err)

		err = repository.UpdateSetCoverByID(entity.ID, "cover1")
		assert.Nil(t, err)

		var updatedEntity catalogEntryEntity
		err = db.Get(&updatedEntity, "select * from catalog_entries where id = $1", entity.ID)
		assert.Nil(t, err)
		assert.Equal(t, "cover1", updatedEntity.Cover.String)
	})
}

func TestCatalogRepository_FindAllBooksWithoutCover(t *testing.T) {
	withDB(func(db *DB) {
		book1, book2, _, _, catalog := generateCatalog(db)

		repository := NewCatalogRepository(db, util.NewUUIDGenerator())
		err := repository.Save(catalog)
		assert.Nil(t, err)

		books, err := repository.FindAllBooksWithoutCover()
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.Equal(t, book1.Name, books[0].Name)
		assert.Equal(t, book2.Name, books[1].Name)
	})
}

func generateCatalog(db *DB) (domain.CatalogEntry, domain.CatalogEntry, domain.CatalogEntry, domain.CatalogEntry, domain.Catalog) {
	book1 := domain.CatalogEntry{Name: "book1", Path: "ebooks/book1", IsDirectory: false, Type: domain.EPUB, Children: nil}
	book2 := domain.CatalogEntry{Name: "book2", Path: "ebooks/book2", IsDirectory: false, Type: domain.EPUB, Children: nil}
	book3 := domain.CatalogEntry{Name: "book3", Path: "ebooks/book3", IsDirectory: false, Type: domain.CBZ, Children: nil}
	ebooks := domain.CatalogEntry{Name: "ebooks", Path: "ebooks", IsDirectory: true, Type: domain.DIRECTORY, Children: []domain.CatalogEntry{book1, book2}}
	catalog := domain.Catalog{Root: ebooks, SourceID: "sourceID"}
	db.MustExec("insert into sources(id, name, path) VALUES ('sourceID', 'sourceName', 'sourcePath')")
	return book1, book2, book3, ebooks, catalog
}
