package domain

import (
	"errors"
	"path/filepath"
)

//go:generate mockgen -destination=../mock/domain/catalog.go -source=catalog.go

type CatalogEntryType string
type CatalogEntryMIMEType string
type CatalogEntryTypeAnalyzer func(path string) CatalogEntryType

const (
	MimeEpub CatalogEntryMIMEType = "application/epub+zip"
	MimeCbz  CatalogEntryMIMEType = "application/vnd.comicbook+zip"
	MimeCbr  CatalogEntryMIMEType = "application/vnd.comicbook+rar"
	MimeAzw3 CatalogEntryMIMEType = "application/vnd.amazon.ebook"
	MimeMobi CatalogEntryMIMEType = "application/x-mobipocket-ebook"

	DIRECTORY CatalogEntryType = ""
	EPUB      CatalogEntryType = "EPUB"
	CBZ       CatalogEntryType = "CBZ"
	CBR       CatalogEntryType = "CBR"
	AZW3      CatalogEntryType = "AZW3"
	MOBI      CatalogEntryType = "MOBI"
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
	Cover       string
}

type CatalogRepository interface {
	Save(catalog Catalog) error
	FindByID(id string) (CatalogEntry, error)
	FindAllRootDirectories() ([]CatalogEntry, error)
	FindAllBooks() ([]CatalogEntry, error)
	FindAllByParentCatalogEntryID(parentCatalogEntryID string) ([]CatalogEntry, error)
	FindAllByParentCatalogEntryIDInPage(parentCatalogEntryID string, page int, pageSize int) ([]CatalogEntry, error)
	FindAllBooksInPage(page int, pageSize int) ([]CatalogEntry, error)
	CountBooks() (int, error)
	CountByParentCatalogEntryID(parentCatalogEntryID string) (int, error)
	UpdateSetCoverByID(id string, cover string) error
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

func GetMimeType(entryType CatalogEntryType) CatalogEntryMIMEType {
	switch entryType {
	case EPUB:
		return MimeEpub
	case CBZ:
		return MimeCbz
	case CBR:
		return MimeCbr
	case AZW3:
		return MimeAzw3
	case MOBI:
		return MimeMobi
	default:
		return ""
	}
}
