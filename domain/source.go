package domain

import "time"

//go:generate mockgen -destination=../mock/domain/source.go -source=source.go

type Source struct {
	ID        string
	Name      string
	Path      string
	UpdatedAt time.Time
}

type SourceRepository interface {
	Insert(source Source) (string, error)
	FindAll() ([]Source, error)
	FindByID(id string) (Source, error)
	DeleteByID(id string) error
}
