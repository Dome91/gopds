package services

import "gopds/domain"

type FetchCatalogEntryByID func(id string) (domain.CatalogEntry, error)
type FetchCatalogRoots func() ([]domain.CatalogEntry, error)
type FetchAllBooksInPage func(page int, pageSize int) ([]domain.CatalogEntry, error)
type CountAllBooks func() (int, error)

func FetchCatalogEntryByIDProvider(repository domain.CatalogRepository) FetchCatalogEntryByID {
	return func(id string) (domain.CatalogEntry, error) {
		return repository.FindByID(id)
	}
}

func FetchCatalogRootsProvider(repository domain.CatalogRepository) FetchCatalogRoots {
	return func() ([]domain.CatalogEntry, error) {
		return repository.FindAllRoots()
	}
}

func FetchAllBooksInPageProvider(repository domain.CatalogRepository) FetchAllBooksInPage {
	return func(page int, pageSize int) ([]domain.CatalogEntry, error) {
		return repository.FindAllBooksInPage(page, pageSize)
	}
}

func CountAllBooksProvider(repository domain.CatalogRepository) CountAllBooks {
	return func() (int, error) {
		return repository.CountAllBooks()
	}
}
