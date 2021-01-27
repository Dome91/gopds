package services

import (
	"gopds/domain"
)

type GenerateOPDSRootFeed func() domain.Feed
type GenerateOPDSAllFeed func() (domain.Feed, error)
type GenerateOPDSFoldersFeed func() (domain.Feed, error)
type GenerateOPDSFeedByID func(id string) (domain.Feed, error)

func GenerateOPDSRootFeedProvider() GenerateOPDSRootFeed {
	return func() domain.Feed {
		allAcquisitionEntry := domain.AllAcquisitionEntry()
		foldersAcquisitionEntry := domain.FoldersAcquisitionEntry()
		return domain.RootFeed(allAcquisitionEntry, foldersAcquisitionEntry)
	}
}

func GenerateOPDSAllFeedProvider(repository domain.CatalogRepository) GenerateOPDSAllFeed {
	return func() (domain.Feed, error) {
		catalogEntries, err := repository.FindAllFiles()
		if err != nil {
			return domain.Feed{}, err
		}

		var entries []domain.Entry
		for _, catalogEntry := range catalogEntries {
			entry := domain.Entry{ID: catalogEntry.ID, Title: catalogEntry.Name, Links: []domain.Link{domain.FileAcquisitionLink(catalogEntry)}}
			entries = append(entries, entry)
		}

		return domain.AllFeed(entries), nil
	}
}

func GenerateOPDSFoldersFeedProvider(repository domain.CatalogRepository) GenerateOPDSFoldersFeed {
	return func() (domain.Feed, error) {
		catalogEntries, err := repository.FindAllRoots()
		if err != nil {
			return domain.Feed{}, err
		}

		var entries []domain.Entry
		for _, catalogEntry := range catalogEntries {
			entry := domain.Entry{ID: catalogEntry.ID, Title: catalogEntry.Name, Links: []domain.Link{domain.DirectoryAcquisitionLink(catalogEntry.ID)}}
			entries = append(entries, entry)
		}

		links := []domain.Link{domain.SelfLink("folders"), domain.StartLink()}
		return domain.NewFeed("Catalog Folders", entries, links), nil
	}
}

func GenerateOPDSFeedByIDProvider(repository domain.CatalogRepository) GenerateOPDSFeedByID {
	return func(id string) (domain.Feed, error) {
		catalogEntry, err := repository.FindByID(id)
		if err != nil {
			return domain.Feed{}, err
		}

		children, err := repository.FindAllByParentCatalogEntryID(catalogEntry.ID)
		if err != nil {
			return domain.Feed{}, err
		}

		var entries []domain.Entry
		for _, child := range children {
			var entry domain.Entry
			if child.IsDirectory {
				entry = domain.Entry{ID: child.ID, Title: child.Name, Links: []domain.Link{domain.DirectoryAcquisitionLink(child.ID)}}
			} else {
				entry = domain.Entry{ID: child.ID, Title: child.Name, Links: []domain.Link{domain.FileAcquisitionLink(child)}}
			}
			entries = append(entries, entry)
		}

		links := []domain.Link{domain.SelfLink(id), domain.StartLink()}
		return domain.NewFeed(catalogEntry.Name, entries, links), nil
	}
}
