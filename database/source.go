package database

import (
	"github.com/google/uuid"
	"gopds/domain"
	"time"
)

type SourceRepository struct {
	db *DB
}

func NewSourceRepository(db *DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func (s *SourceRepository) Insert(source domain.Source) (string, error) {
	id := uuid.New().String()
	_, err := s.db.Exec("insert into sources(id, name, path) values($1, $2, $3)", id, source.Name, source.Path)
	return id, err
}

func (s *SourceRepository) FindAll() ([]domain.Source, error) {
	var entities []sourceEntity
	err := s.db.Select(&entities, "select * from sources")
	return s.mapAllToDomain(entities), err
}

func (s *SourceRepository) FindByID(id string) (domain.Source, error) {
	var entity sourceEntity
	err := s.db.Get(&entity, "select * from sources where id = $1", id)
	return s.mapToDomain(entity), err
}

func (s *SourceRepository) DeleteByID(id string) error {
	_, err := s.db.Exec("delete from sources where id  = $1", id)
	return err
}

func (s *SourceRepository) mapAllToDomain(entities []sourceEntity) []domain.Source {
	sources := make([]domain.Source, len(entities))
	for index, entity := range entities {
		sources[index] = s.mapToDomain(entity)
	}

	return sources
}

func (s *SourceRepository) mapToDomain(entity sourceEntity) domain.Source {
	return domain.Source{
		ID:        entity.ID,
		Name:      entity.Name,
		Path:      entity.Path,
		UpdatedAt: entity.UpdatedAt,
	}
}

type sourceEntity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
