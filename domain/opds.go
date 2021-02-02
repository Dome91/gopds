package domain

import (
	"encoding/xml"
	"fmt"
)

const (
	acquisitionType = "application/atom+xml;profile=opds-catalog;kind=acquisition"
	allTitle        = "All Catalog Entries"
	id              = "GOPDS"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
	Links   []Link   `xml:"link"`
}

func NewFeed(title string, entries []Entry, links []Link) Feed {
	return Feed{
		ID:      id,
		Title:   title,
		Entries: entries,
		Links:   links,
	}
}

func RootFeed(entries ...Entry) Feed {
	return Feed{
		ID:      id,
		Title:   "Root",
		Entries: entries,
		Links:   []Link{StartLink(), SelfLink("")},
	}
}

func AllFeed(entries []Entry) Feed {
	return Feed{
		ID:      id,
		Title:   allTitle,
		Entries: entries,
		Links:   []Link{StartLink(), SelfLink("all")},
	}
}

type Link struct {
	XMLName xml.Name             `xml:"link"`
	Href    string               `xml:"href,attr"`
	Rel     string               `xml:"rel,attr"`
	Type    CatalogEntryMIMEType `xml:"type,attr"`
}

func SelfLink(href string) Link {
	return Link{
		Href: "/opds/" + href,
		Rel:  "self",
		Type: acquisitionType,
	}
}

func StartLink() Link {
	return Link{
		Href: "/opds",
		Rel:  "start",
		Type: acquisitionType,
	}
}

func DirectoryAcquisitionLink(id string) Link {
	return Link{
		Href: "/opds/" + id,
		Rel:  "http://opds-spec.org/crawlable",
		Type: acquisitionType,
	}
}

func BookAcquisitionLink(catalogEntry CatalogEntry) Link {
	href := fmt.Sprintf("/opds/%s/download", catalogEntry.ID)
	return Link{
		Href: href,
		Rel:  "http://opds-spec.org/acquisition/",
		Type: GetMimeType(catalogEntry.Type),
	}
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Links   []Link   `xml:"link"`
}

func AllAcquisitionEntry() Entry {
	return Entry{
		ID:    "all",
		Title: allTitle,
		Links: []Link{DirectoryAcquisitionLink("all")},
	}
}

func DirectoriesAcquisitionEntry() Entry {
	return Entry{
		ID:    "directories",
		Title: "Catalog Directories",
		Links: []Link{DirectoryAcquisitionLink("directories")},
	}
}
