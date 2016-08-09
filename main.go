package main

import (
	"flag"
	"fmt"
	"log"
)

// http://technosophos.com/2014/06/11/compile-time-string-in-go.html
// go build -ldflags "-x main.version=$(git describe --tags)"
var version string

var verbose = flag.Bool("v", false, "verbose output")

func main() {
	var (
		bindAddr    = flag.String("addr", ":8080", "http listen address")
		indexPath   = flag.String("indexPath", "indexes/search.bleve", "path of the bleve index")
		hugoPath    = flag.String("hugoPath", ".", "path of the hugo site")
		showVersion = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()
	if flag.Arg(0) == "version" || *showVersion {
		fmt.Println("hugo-search version", version)
		return
	}
	if flag.NArg() > 0 {
		flag.PrintDefaults()
	}
	log.SetFlags(0)

	buildIndexFromSite(*hugoPath, *indexPath)
	startSearchServer(*bindAddr, *indexPath)
}

func checkFatal(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
