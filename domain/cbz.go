package domain

import (
	"archive/zip"
	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"sort"
	"strings"
)

func ForEveryFileInCBZDo(entry CatalogEntry, predicate CatalogEntryFilePredicate, consumer func(file archiver.File) error) error {
	walker, err := getSupportedWalker(entry)
	if err != nil {
		return err
	}

	return walker.Walk(entry.Path, func(f archiver.File) error {
		if predicate(f) {
			return consumer(f)
		}

		return nil
	})
}

func GetFilePathsInCatalogEntryInAlphabeticalOrder(entry CatalogEntry, predicate CatalogEntryFilePredicate) ([]string, error) {
	walker, err := getSupportedWalker(entry)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	err = walker.Walk(entry.Path, func(f archiver.File) error {
		if predicate(f) {
			filePath, err := GetFilePath(entry.Type, f)
			if err != nil {
				log.Error(err)
				return nil
			}

			filePaths = append(filePaths, filePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(filePaths, func(i, j int) bool { return strings.ToLower(filePaths[i]) < strings.ToLower(filePaths[j]) })
	return filePaths, nil
}

func GetFilePath(entryType CatalogEntryType, f archiver.File) (string, error) {
	switch entryType {
	case CBZ:
		header := f.Header.(zip.FileHeader)
		return header.Name, nil
	default:
		return "", ErrUnsupportedFiletype
	}
}

type CatalogEntryFilePredicate func(file archiver.File) bool

func IsImage(file archiver.File) bool {
	extension := filepath.Ext(file.FileInfo.Name())
	return extension == ".jpg" || extension == ".jpeg" || extension == ".png"
}

func IsDirectory(file archiver.File) bool {
	return file.IsDir()
}

func Not(predicate CatalogEntryFilePredicate) CatalogEntryFilePredicate {
	return func(file archiver.File) bool {
		return !predicate(file)
	}
}

func And(p1 CatalogEntryFilePredicate, p2 CatalogEntryFilePredicate) CatalogEntryFilePredicate {
	return func(file archiver.File) bool {
		return p1(file) && p2(file)
	}
}

func OnlyImages(file archiver.File) bool {
	return And(Not(IsDirectory), IsImage)(file)
}

func getSupportedWalker(entry CatalogEntry) (archiver.Walker, error) {
	switch entry.Type {
	case CBZ:
		return archiver.NewZip(), nil
	default:
		return nil, ErrUnsupportedFiletype
	}
}
