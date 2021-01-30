package domain

import (
	"errors"
	"path/filepath"
)

//go:generate mockgen -destination=../mock/domain/catalog.go -source=catalog.go

type CatalogEntryType string
type CatalogEntryTypeAnalyzer func(path string) CatalogEntryType

const (
	DIRECTORY CatalogEntryType = ""
	EPUB      CatalogEntryType = "application/epub+zip"
	CBZ       CatalogEntryType = "application/vnd.comicbook+zip"
	CBR       CatalogEntryType = "application/vnd.comicbook+rar"
	AZW3      CatalogEntryType = "application/vnd.amazon.ebook"
	MOBI      CatalogEntryType = "application/x-mobipocket-ebook"
)

var (
	ErrUnsupportedFiletype = errors.New("unsupported filetype")
	analyzers              = [...]CatalogEntryTypeAnalyzer{
		catalogEntryTypeAnalyzerProvider(".cbz", CBZ),
		catalogEntryTypeAnalyzerProvider(".cbr", CBR),
		catalogEntryTypeAnalyzerProvider(".epub", EPUB),
		catalogEntryTypeAnalyzerProvider(".mobi", MOBI),
		catalogEntryTypeAnalyzerProvider(".azw3", AZW3),
	}
)

type Catalog struct {
	ID       string
	Root     CatalogEntry
	SourceID string
}

type CatalogEntry struct {
	ID          string
	Name        string
	Path        string
	IsDirectory bool
	Children    []CatalogEntry
	Type        CatalogEntryType
}

type CatalogRepository interface {
	Save(catalog Catalog) error
	FindByID(id string) (CatalogEntry, error)
	FindAllRoots() ([]CatalogEntry, error)
	FindAllBooks() ([]CatalogEntry, error)
	FindAllByParentCatalogEntryID(parentCatalogEntryID string) ([]CatalogEntry, error)
	FindAllBooksInPage(page int, pageSize int) ([]CatalogEntry, error)
	CountAllBooks() (int, error)
}

func catalogEntryTypeAnalyzerProvider(extension string, entryType CatalogEntryType) CatalogEntryTypeAnalyzer {
	return func(path string) CatalogEntryType {
		ext := filepath.Ext(path)
		if ext == extension {
			return entryType
		}

		return ""
	}
}

func DetermineCatalogEntryType(path string) (CatalogEntryType, error) {
	for _, analyzer := range analyzers {
		entryType := analyzer(path)
		if entryType != "" {
			return entryType, nil
		}
	}

	return "", ErrUnsupportedFiletype
}
