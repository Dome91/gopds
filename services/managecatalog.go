package services

import "gopds/domain"

type FetchCatalogEntryByID func (id string) (domain.CatalogEntry, error)

func FetchCatalogEntryByIDProvider(repository domain.CatalogRepository) FetchCatalogEntryByID {
	return func(id string) (domain.CatalogEntry, error) {
		return repository.FindByID(id)
	}
}
