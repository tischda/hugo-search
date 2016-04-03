package main

import (
	"testing"

	"github.com/blevesearch/bleve"
)

const testIndexPath = "test/indexes/search.bleve"

// checks the actual index creation and validity
func TestBuildIndex(t *testing.T) {
	buildIndexFromSite(testHugoPath, testIndexPath)
	index := openIndex(t, testIndexPath)
	defer index.Close()
	queryIndex(t, index)
}

// loads the index
func openIndex(t *testing.T, path string) bleve.Index {
	index, err := bleve.OpenUsing(path, map[string]interface{}{"read_only": true})
	if err != nil {
		t.Errorf("error opening index %s: %v", path, err)
	} else if index == nil {
		t.Error("null index")
	}
	return index
}

// queries the index
func queryIndex(t *testing.T, index bleve.Index) {
	query := bleve.NewTermQuery("lorem")
	request := bleve.NewSearchRequest(query)
	result, err := index.Search(request)
	if err != nil {
		t.Error(err)
	}
	if result.Total < 1 {
		t.Error("No hits for 'lorem', expected at least one.")
	}
}
