package main

import (
	"log"
	"net/http"

	"path"

	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"
	"github.com/rs/cors"
)

func startSearchServer(addr string, indexPath string) {
	indexName := path.Base(indexPath)
	index := registerIndex(indexPath, indexName)
	defer unregisterIndex(index, indexName)
	handler := getCorsHandler(indexName)

	log.Printf("Search server listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// registers the index by its name so that handler can use it
func registerIndex(indexPath string, indexName string) bleve.Index {
	if *verbose {
		log.Printf("Registering index: %s", indexPath)
	}
	index, err := bleve.OpenUsing(indexPath, map[string]interface{}{"read_only": true})
	checkFatal(err)
	bleveHttp.RegisterIndexName(indexName, index)
	return index
}

// unregister and close the index
func unregisterIndex(index bleve.Index, indexName string) {
	bleveHttp.UnregisterIndexByName(indexName)
	index.Close()
}

//  Cross Origin Resource Sharing (https://www.w3.org/TR/cors/)
func getCorsHandler(indexName string) http.Handler {

	// list of indexes
	mux := http.NewServeMux()
	mux.HandleFunc("/api", bleveHttp.NewListIndexesHandler().ServeHTTP)

	// actual search handler
	searchHandler := bleveHttp.NewSearchHandler(indexName)
	mux.HandleFunc("/api/"+indexName+"/_search", searchHandler.ServeHTTP)
	return cors.Default().Handler(mux)
}
