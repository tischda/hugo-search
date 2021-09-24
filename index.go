package main

import (
	"log"

	"os"

	"github.com/blevesearch/bleve"
	"github.com/gohugoio/hugo/resources/page"
)

// builds the search index by passing all pages of hugo site that have a title to the indexer
func buildIndexFromSite(theHugoPath string, theIndexPath string) {
	pages := readSitePages(theHugoPath)
	index := createIndex(theIndexPath)
	defer index.Close()
	for _, page := range pages {
		if pageHasTitle(page) && page.Type() != "search" {
			addPageToIndex(index, page)
		}
	}
}

// creates the index from scratch (does not reuse existing index)
func createIndex(path string) bleve.Index {
	if *verbose {
		log.Println("Creating Index:", path)
	}

	// index_meta.go, line 59: os.Mkdir(path, 0700) fails if parent directory missing
	err := os.MkdirAll(path, 0700)
	exitOnError(err)

	// always recreate full index (otherwise search returns deleted pages)
	err = os.RemoveAll(path)
	exitOnError(err)

	index, err := bleve.New(path, bleve.NewIndexMapping())
	exitOnError(err)
	return index
}

// adds a hugo page to the bleve search index
func addPageToIndex(index bleve.Index, p page.Page) {
	link := p.RelPermalink()
	exitOnError(index.Index(link, newIndexEntry(p)))
	if *verbose {
		log.Printf("Indexed: %s [%s]", p.File().Path(), p.Title())
	}
}
