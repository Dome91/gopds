package domain

import (
	"github.com/mholt/archiver/v3"
	"path/filepath"
	"sort"
	"strings"
)

func ForEveryFileInCatalogEntryDo(entry CatalogEntry, predicate CatalogEntryFilePredicate, consumer func(file archiver.File) error) error {
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

func GetFilenamesInCatalogEntryInAlphabeticalOrder(entry CatalogEntry, predicate CatalogEntryFilePredicate) ([]string, error) {
	walker, err := getSupportedWalker(entry)
	if err != nil {
		return nil, err
	}

	var filenames []string
	err = walker.Walk(entry.Path, func(f archiver.File) error {
		if predicate(f) {
			filenames = append(filenames, f.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(filenames, func(i, j int) bool { return strings.ToLower(filenames[i]) < strings.ToLower(filenames[j]) })
	return filenames, nil
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
