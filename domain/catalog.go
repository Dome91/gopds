package domain

//go:generate mockgen -destination=../mock/domain/catalog.go -source=catalog.go

type CatalogEntryType string

const (
	DIRECTORY CatalogEntryType = ""
	EPUB      CatalogEntryType = "application/epub+zip"
	CBZ       CatalogEntryType = "application/vnd.comicbook+zip"
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
	FindAllFiles() ([]CatalogEntry, error)
	FindAllByParentCatalogEntryID(parentCatalogEntryID string) ([]CatalogEntry, error)
}
