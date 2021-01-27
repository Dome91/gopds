package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopds/domain"
	"testing"
)

func TestSourceRepository_Insert(t *testing.T) {
	withDB(func(db *DB) {
		source := domain.Source{Name: "name", Path: "path"}

		_, err := NewSourceRepository(db).Insert(source)
		assert.Nil(t, err)

		var entities []sourceEntity
		err = db.Select(&entities, "select * from sources")
		assert.Nil(t, err)
		assert.Len(t, entities, 1)

		entity := entities[0]
		assert.Equal(t, source.Name, entity.Name)
		assert.Equal(t, source.Path, entity.Path)
		assert.NotPanics(t, func() {
			uuid.MustParse(entity.ID)
		})
		assert.NotNil(t, entity.UpdatedAt)
	})
}

func TestSourceRepository_FindAll(t *testing.T) {
	withDB(func(db *DB) {
		source1 := domain.Source{Name: "name1", Path: "path1"}
		source2 := domain.Source{Name: "name2", Path: "path2"}

		repository := NewSourceRepository(db)
		_, err := repository.Insert(source1)
		assert.Nil(t, err)
		_, err = repository.Insert(source2)
		assert.Nil(t, err)

		sources, err := repository.FindAll()
		assert.Nil(t, err)
		assert.Len(t, sources, 2)
		assert.Equal(t, source1.Name, sources[0].Name)
		assert.Equal(t, source1.Path, sources[0].Path)
		assert.Equal(t, source2.Name, sources[1].Name)
		assert.Equal(t, source2.Path, sources[1].Path)
	})
}

func TestSourceRepository_FindByID(t *testing.T) {
	withDB(func(db *DB) {
		repository := NewSourceRepository(db)

		source := domain.Source{Name: "name", Path: "path"}
		id, err := repository.Insert(source)
		assert.Nil(t, err)

		foundSource, err := repository.FindByID(id)
		assert.Nil(t, err)
		assert.Equal(t, source.Name, foundSource.Name)
		assert.Equal(t, source.Path, foundSource.Path)
		assert.NotNil(t, foundSource.UpdatedAt)
	})
}

func TestSourceRepository_DeleteByID(t *testing.T) {
	withDB(func(db *DB) {
		_, err := db.Exec("insert into sources (id, name, path) values ('id1', 'name', 'path')")
		require.Nil(t, err)
		_, err = db.Exec("insert into catalog_entries (id, name, path, is_directory, type, parent_catalog_entry, source) VALUES ('id2', 'name2', 'path2',  true, 'type1', null, 'id1')")
		require.Nil(t, err)
		_, err = db.Exec("insert into catalog_entries (id, name, path, is_directory, type,  parent_catalog_entry, source) VALUES ('id3', 'name3', 'path3',  true, 'type1', 'id2', 'id1')")
		require.Nil(t, err)

		err = NewSourceRepository(db).DeleteByID("id1")
		assert.Nil(t, err)

		var count int
		err = db.Get(&count, "select count(*) from sources")
		assert.Nil(t, err)
		assert.Zero(t, count)

		err = db.Get(&count, "select count(*) from catalog_entries")
		assert.Nil(t, err)
		assert.Zero(t, count)
	})
}
