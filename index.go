package main

import (
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/gohugoio/hugo/resources/page"
)

// builds the search index by passing all pages of hugo site that have a title to the indexer
func buildIndexFromSite(cfg *Config) {
	pages := readSitePages(cfg.hugoPath)
	index := createIndex(cfg.indexPath, cfg.verbose)
	defer func() {
		if err := index.Close(); err != nil {
			log.Printf("Error closing index %s: %v", cfg.indexPath, err)
		}
	}()
	for _, page := range pages {
		if pageHasTitle(page, false) && page.Type() != "search" {
			addPageToIndex(index, page, cfg.verbose)
		}
	}
}

// creates the index from scratch (does not reuse existing index)
func createIndex(path string, verbose bool) bleve.Index {
	if verbose {
		log.Println("Creating Index:", path)
	}

	// index_meta.go, line 59: os.Mkdir(path, 0700) fails if parent directory missing
	err := os.MkdirAll(path, 0700)
	exitOnError(err)

	// always recreate full index (otherwise search returns deleted pages)
	err = os.RemoveAll(path)
	exitOnError(err)

	index, err := bleve.New(path, mapping.NewIndexMapping())
	exitOnError(err)
	return index
}

func pageHasTitle(p page.Page, verbose bool) bool {
	found := len(p.Title()) > 0
	if !found && verbose {
		log.Println("WARN: Title is missing in file metadata:", p.File().Path())
	}
	return found
}

// adds a hugo page to the bleve search index
func addPageToIndex(index bleve.Index, p page.Page, verbose bool) {
	link := p.RelPermalink()
	exitOnError(index.Index(link, newIndexEntry(p)))
	if verbose {
		log.Printf("Indexed: %s [%s]", p.File().Path(), p.Title())
	}
}
