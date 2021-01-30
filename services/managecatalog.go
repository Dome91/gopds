package services

import "gopds/domain"

type FetchCatalogEntryByID func(id string) (domain.CatalogEntry, error)
type FetchCatalogEntriesByParentCatalogEntryIDInPage func(parentCatalogEntryID string, page int, pageSize int) ([]domain.CatalogEntry, error)
type FetchCatalogRootDirectories func() ([]domain.CatalogEntry, error)
type FetchAllBooksInPage func(page int, pageSize int) ([]domain.CatalogEntry, error)
type CountBooks func() (int, error)
type CountCatalogEntriesByParentCatalogEntryID func(parentCatalogEntryID string) (int, error)

func FetchCatalogEntryByIDProvider(repository domain.CatalogRepository) FetchCatalogEntryByID {
	return func(id string) (domain.CatalogEntry, error) {
		return repository.FindByID(id)
	}
}

func FetchCatalogEntriesByParentCatalogEntryIDInPageProvider(repository domain.CatalogRepository) FetchCatalogEntriesByParentCatalogEntryIDInPage {
	return func(parentCatalogEntryID string, page int, pageSize int) ([]domain.CatalogEntry, error) {
		return repository.FindAllByParentCatalogEntryIDInPage(parentCatalogEntryID, page, pageSize)
	}
}

func FetchCatalogRootDirectoriesProvider(repository domain.CatalogRepository) FetchCatalogRootDirectories {
	return func() ([]domain.CatalogEntry, error) {
		return repository.FindAllRootDirectories()
	}
}

func FetchAllBooksInPageProvider(repository domain.CatalogRepository) FetchAllBooksInPage {
	return func(page int, pageSize int) ([]domain.CatalogEntry, error) {
		return repository.FindAllBooksInPage(page, pageSize)
	}
}

func CountBooksProvider(repository domain.CatalogRepository) CountBooks {
	return func() (int, error) {
		return repository.CountBooks()
	}
}

func CountCatalogEntriesByParentCatalogEntryIDProvider(repository domain.CatalogRepository) CountCatalogEntriesByParentCatalogEntryID {
	return func(parentCatalogEntryID string) (int, error) {
		return repository.CountByParentCatalogEntryID(parentCatalogEntryID)
	}
}
