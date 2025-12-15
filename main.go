package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// https://goreleaser.com/cookbooks/using-main.version/
var (
	name    string
	version string
	date    string
	commit  string
)

// flags
type Config struct {
	bindAddr  string
	hugoPath  string
	indexPath string
	verbose   bool
	help      bool
	version   bool
}

func initFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.bindAddr, "a", ":8080", "")
	flag.StringVar(&cfg.bindAddr, "bindAddr", ":8080", "http listen address")
	flag.StringVar(&cfg.hugoPath, "h", ".", "")
	flag.StringVar(&cfg.hugoPath, "hugoPath", ".", "path of the hugo site")
	flag.StringVar(&cfg.indexPath, "i", "indexes/search.bleve", "")
	flag.StringVar(&cfg.indexPath, "indexPath", "indexes/search.bleve", "path of the bleve index")
	flag.BoolVar(&cfg.verbose, "vv", false, "")
	flag.BoolVar(&cfg.verbose, "verbose", false, "verbose output")
	flag.BoolVar(&cfg.help, "?", false, "")
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	return cfg
}

func main() {
	log.SetFlags(0)
	cfg := initFlags()
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+name+` [OPTIONS]

Starts a search engine for a Hugo site using a Bleve index.

OPTIONS:

  -a, --bindAddr"
          http listen address (default ":8080")
  -h, --hugoPath
          path to the hugo site (default ".")
  -i, --indexPath
          path to the bleve index (default "indexes/search.bleve")
  -vv, --verbose
          verbose output
  -?, --help
          display this help message
  -v, --version
          print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  $ "+name+` --hugoPath test
  All pages: 12, regular pages: 5
  Search server listening on :8080
  `)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || cfg.version {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", name, version, date, commit)
		return
	}

	if cfg.help {
		flag.Usage()
		return
	}

	if !flag.Parsed() || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	log.SetFlags(0)

	buildIndexFromSite(cfg)
	startSearchServer(cfg)
}

func exitOnError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
