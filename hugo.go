package main

import (
	"log"

	"path/filepath"

	"os"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
)

// Returns all pages of the hugo site located at path.
// This function duplicates code from InitializeConfig() (cf. commands/hugo.go)
// because it does not allow us to set the workdingDir
func readSitePages(path string) hugolib.Pages {

	var dir string
	if path != "" {
		dir, _ = filepath.Abs(path)

	} else {
		dir, _ = os.Getwd()
	}
	osFs := hugofs.Os

	// WorkdingDir is not evaluated here
	cfg := hugolib.ConfigSourceDescriptor{Fs: osFs, Path: dir}
	config, _, err := hugolib.LoadConfig(cfg)
	checkFatal(err)

	// We still need to set workdingDir
	config.Set("workingDir", dir)

	fs := hugofs.NewFrom(osFs, config)

	// cf. hugolib/hugo_sites_build_test.go
	sites, err := hugolib.NewHugoSites(deps.DepsCfg{Cfg: config, Fs: fs})
	checkFatal(err)

	err = sites.Build(hugolib.BuildCfg{SkipRender: true})
	checkFatal(err)

	return sites.Pages()
}

// checks if the page has a title (that will appear in the search result)
func pageHasTitle(page *hugolib.Page) (foundTitle bool) {
	foundTitle = len(page.Title()) > 0
	if !foundTitle && *verbose {
		log.Println("WARN: Title metadata missing in document:", page.File.Path())
	}
	return
}

// checks if the page content is valid to be indexed
//
// 'title-page-3:page'		yes
// ':page'			yes	ok page has no title, but that could be dealt with elsewhere
// 'title-page-1:page'		yes
// 'title-page-2:page'		yes
// 'Search Results:page'	no	dynamic content, do not index (wish there was a kind 'searchResults')
// 'Categories:taxonomyTerm'    no	dynamic content, do not index
// 'Fails:section'		yes
// 'Folder1s:section'		yes
// 'Tag1:taxonomy'		no	dynamic content, do not index
// 'Tag2:taxonomy'		no	dynamic content, do not index
// 'Tags:taxonomyTerm'		no	dynamic content, do not index
// 'hugo-search test site:home'	yes
func pageHasValidContent(page *hugolib.Page) bool {
	switch page.Kind {
	case "page":
		if page.Title() == "Search Results" {
			break
		}
		fallthrough
	case "section":
		fallthrough
	case "home":
		return true
	}
	if *verbose {
		log.Printf("Ignoring: %s [%s] of kind: %s", page.File.Path(), page.Title(), page.Kind)
	}
	return false
}
