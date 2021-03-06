package services

import (
	log "github.com/sirupsen/logrus"
	"gopds/domain"
	"gopds/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SynchronizeCatalog func(sourceID string) error

func SynchronizeCatalogProvider(sourceRepository domain.SourceRepository, catalogRepository domain.CatalogRepository, bus util.Bus) SynchronizeCatalog {
	return func(sourceID string) error {
		source, err := sourceRepository.FindByID(sourceID)
		if err != nil {
			return err
		}

		root, err := processCatalogEntry(source.Path)
		if err != nil {
			return err
		}

		catalog := domain.Catalog{Root: root, SourceID: sourceID}
		err = catalogRepository.Save(catalog)
		if err != nil {
			return err
		}

		booksWithoutCover, err := catalogRepository.FindAllBooksWithoutCover()
		if err != nil {
			return err
		}
		domain.SendGenerateCoverEvents(booksWithoutCover, bus)
		return nil
	}
}

func processCatalogEntry(path string) (domain.CatalogEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return domain.CatalogEntry{}, err
	}

	if info.IsDir() {
		return processDirectoryCatalogEntry(info, path)
	}

	return processFileCatalogEntry(info, path)
}

func processDirectoryCatalogEntry(info os.FileInfo, path string) (domain.CatalogEntry, error) {
	subDirectoryInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return domain.CatalogEntry{}, err
	}

	var children []domain.CatalogEntry
	for _, subDirectoryInfo := range subDirectoryInfos {
		subDirectoryPath, err := filepath.Abs(filepath.Join(path, subDirectoryInfo.Name()))
		if err != nil {
			log.Errorf("failed to get absolute path: %s", err.Error())
			continue
		}

		childEntry, err := processCatalogEntry(subDirectoryPath)
		if err != nil {
			log.Errorf("failed to process catalog entry %s: %s", subDirectoryInfo.Name(), err.Error())
			continue
		}

		children = append(children, childEntry)
	}

	return domain.CatalogEntry{Name: info.Name(), Path: path, IsDirectory: true, Children: children, Type: domain.DIRECTORY}, nil
}

func processFileCatalogEntry(info os.FileInfo, path string) (domain.CatalogEntry, error) {
	catalogEntryType, err := domain.DetermineCatalogEntryType(path)
	if err != nil {
		return domain.CatalogEntry{}, err
	}

	return domain.CatalogEntry{Name: info.Name(), Path: path, IsDirectory: false, Children: nil, Type: catalogEntryType}, nil
}
