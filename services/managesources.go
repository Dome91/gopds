package services

import (
	"errors"
	"gopds/domain"
	"os"
	"path/filepath"
)

var ErrInvalidPath = errors.New("path is not valid")

type CreateSource func(name string, path string) error
type FetchAllSources func() ([]domain.Source, error)
type DeleteSource func(id string) error

func CreateSourceProvider(repository domain.SourceRepository) CreateSource {
	validatePath := func(path string) error {
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			return ErrInvalidPath
		}

		if !fileInfo.IsDir() {
			return ErrInvalidPath
		}

		return nil
	}

	return func(name string, path string) error {
		err := validatePath(path)
		if err != nil {
			return err
		}

		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		source := domain.Source{Name: name, Path: absolutePath}
		_, err = repository.Insert(source)
		return err
	}
}

func FetchAllSourcesProvider(repository domain.SourceRepository) FetchAllSources {
	return func() ([]domain.Source, error) {
		return repository.FindAll()
	}
}

func DeleteSourceProvider(repository domain.SourceRepository) DeleteSource {
	return func(id string) error {
		return repository.DeleteByID(id)
	}
}
