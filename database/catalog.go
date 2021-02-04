package database

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"gopds/domain"
	"time"
)

type CatalogRepository struct {
	db *DB
	b  domain.Bus
}

func NewCatalogRepository(db *DB, b domain.Bus) *CatalogRepository {
	return &CatalogRepository{db: db, b: b}
}

func (c *CatalogRepository) Save(catalog domain.Catalog) error {
	saver := catalogSaver{db: c.db}
	err := saver.Save(catalog)
	if err != nil {
		return err
	}

	c.sendGenerateCoverEvents(saver.newCatalogEntryIDs)
	return nil
}

func (c *CatalogRepository) sendGenerateCoverEvents(ids []string) {
	for _, id := range ids {
		_, err := c.b.Emit(context.Background(), domain.GenerateCoverTopic, domain.GenerateCoverEvent{ID: id})
		if err != nil {
			log.Error(err)
		}
	}
}

func (c *CatalogRepository) FindByID(id string) (domain.CatalogEntry, error) {
	var entity catalogEntryEntity
	err := c.db.Get(&entity, "select * from catalog_entries where id = $1", id)
	return c.mapToDomain(entity), err
}

func (c *CatalogRepository) FindAllByParentCatalogEntryID(parentCatalogEntryID string) ([]domain.CatalogEntry, error) {
	var entities []catalogEntryEntity
	err := c.db.Select(&entities, "select * from catalog_entries where parent_catalog_entry = $1 order by name", parentCatalogEntryID)
	return c.mapAllToDomain(entities), err
}

func (c *CatalogRepository) FindAllByParentCatalogEntryIDInPage(parentCatalogEntryID string, page int, pageSize int) ([]domain.CatalogEntry, error) {
	var entities []catalogEntryEntity
	err := c.db.Select(&entities, "select * from catalog_entries where parent_catalog_entry = $1 order by name limit $2 offset $3", parentCatalogEntryID, pageSize, page*pageSize)
	return c.mapAllToDomain(entities), err
}

func (c *CatalogRepository) FindAllRootDirectories() ([]domain.CatalogEntry, error) {
	var entities []catalogEntryEntity
	err := c.db.Select(&entities, "select * from catalog_entries where parent_catalog_entry is null order by name")
	return c.mapAllToDomain(entities), err
}

func (c *CatalogRepository) FindAllBooks() ([]domain.CatalogEntry, error) {
	var entities []catalogEntryEntity
	err := c.db.Select(&entities, "select * from catalog_entries where is_directory = false order by name")
	return c.mapAllToDomain(entities), err
}

func (c *CatalogRepository) FindAllBooksInPage(page int, pageSize int) ([]domain.CatalogEntry, error) {
	var entities []catalogEntryEntity
	err := c.db.Select(&entities, "select * from catalog_entries where is_directory = false order by name limit $1 offset $2", pageSize, page*pageSize)
	return c.mapAllToDomain(entities), err
}

func (c *CatalogRepository) CountBooks() (int, error) {
	var count int
	err := c.db.Get(&count, "select count(*) from catalog_entries where is_directory = false")
	return count, err
}

func (c *CatalogRepository) CountByParentCatalogEntryID(parentCatalogEntryID string) (int, error) {
	var count int
	err := c.db.Get(&count, "select count(*) from catalog_entries where parent_catalog_entry = $1", parentCatalogEntryID)
	return count, err
}

func (c *CatalogRepository) UpdateSetCoverByID(id string, cover string) error {
	_, err := c.db.Exec("update catalog_entries set cover = $1 where id = $2", cover, id)
	return err
}

func (c *CatalogRepository) mapToDomain(entity catalogEntryEntity) domain.CatalogEntry {
	return domain.CatalogEntry{
		ID:          entity.ID,
		Name:        entity.Name,
		Path:        entity.Path,
		IsDirectory: entity.IsDirectory,
		Cover:       entity.Cover.String,
		Type:        entity.Type,
	}
}

func (c *CatalogRepository) mapAllToDomain(entities []catalogEntryEntity) []domain.CatalogEntry {
	entries := make([]domain.CatalogEntry, len(entities))
	for index, entity := range entities {
		entries[index] = c.mapToDomain(entity)
	}

	return entries
}

type catalogSaver struct {
	db                 *DB
	newCatalogEntryIDs []string
}

func (c *catalogSaver) Save(catalog domain.Catalog) error {
	return c.db.inTransaction(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("update catalog_entries set found_during_last_sync = false")
		if err != nil {
			return err
		}

		c.saveCatalogEntry(tx, catalog.Root, sql.NullString{Valid: false}, catalog.SourceID)
		_, err = tx.Exec("delete from catalog_entries where found_during_last_sync = false")
		return err
	})
}

func (c *catalogSaver) saveCatalogEntry(tx *sqlx.Tx, entry domain.CatalogEntry, parentID sql.NullString, sourceID string) {
	var entity catalogEntryEntity
	err := tx.Get(&entity, "select * from catalog_entries where path = $1", entry.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			c.insertCatalogEntry(tx, entry, parentID, sourceID)
		} else {
			log.Errorf("get catalog entry by path failed: %s", err.Error())
			return
		}
	} else {
		c.updateCatalogEntry(tx, entity, entry.Children)
	}
}

func (c *catalogSaver) updateCatalogEntry(tx *sqlx.Tx, entity catalogEntryEntity, children []domain.CatalogEntry) {
	_, err := tx.Exec("update catalog_entries set found_during_last_sync = true where path = $1", entity.Path)
	if err != nil {
		log.Errorf("update found_during_last_sync to true for path %s failed: %s", entity.Path, err.Error())
		return
	}

	for _, child := range children {
		c.saveCatalogEntry(tx, child, sql.NullString{String: entity.ID, Valid: true}, entity.SourceID)
	}
}

func (c *catalogSaver) insertCatalogEntry(tx *sqlx.Tx, entry domain.CatalogEntry, parentID sql.NullString, sourceID string) {
	id := uuid.New().String()
	_, err := tx.Exec(
		"insert into catalog_entries(id, name, path, is_directory, type, parent_catalog_entry, source) values ($1, $2, $3, $4, $5, $6, $7)",
		id, entry.Name, entry.Path, entry.IsDirectory, entry.Type, parentID, sourceID)

	if err != nil {
		log.Errorf("inserting new catalog entry failed: %s", err.Error())
		return
	}

	log.Infof("synced %s", entry.Name)
	if !entry.IsDirectory {
		c.newCatalogEntryIDs = append(c.newCatalogEntryIDs, id)
	}

	for _, child := range entry.Children {
		c.saveCatalogEntry(tx, child, sql.NullString{String: id, Valid: true}, sourceID)
	}
}

type catalogEntryEntity struct {
	ID                   string                  `db:"id"`
	Name                 string                  `db:"name"`
	Path                 string                  `db:"path"`
	IsDirectory          bool                    `db:"is_directory"`
	Type                 domain.CatalogEntryType `db:"type"`
	Cover                sql.NullString          `db:"cover"`
	FoundDuringLastSync  bool                    `db:"found_during_last_sync"`
	CreatedAt            time.Time               `db:"created_at"`
	ParentCatalogEntryID sql.NullString          `db:"parent_catalog_entry"`
	SourceID             string                  `db:"source"`
}
