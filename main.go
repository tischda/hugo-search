package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version string

var verbose = flag.Bool("verbose", false, "verbose output")

func main() {
	var (
		bindAddr    = flag.String("addr", ":8080", "http listen address")
		hugoPath    = flag.String("hugoPath", ".", "path of the hugo site")
		indexPath   = flag.String("indexPath", "indexes/search.bleve", "path of the bleve index")
		showVersion = flag.Bool("version", false, "print version and exit")
	)
	// hugo imports the "testing" package so calling flag.PrintDefaults()
	// would display a whole bunch of test.* flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [OPTIONS]\n\nOPTIONS:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -addr <string>\thttp listen address (default \"%s\")\n"+
			"  -hugoPath <string>\tpath of the hugo site (default \"%s\")\n"+
			"  -indexPath <string>\tpath of the bleve index (default \"%s\")\n"+
			"  -verbose\t\tverbose output\n"+
			"  -version\t\tprint version and exit\n", *bindAddr, *hugoPath, *indexPath)
	}
	flag.Parse()
	if !flag.Parsed() || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	if flag.Arg(0) == "version" || *showVersion {
		fmt.Println("hugo-search version", version)
		return
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
