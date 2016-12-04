package main

import (
	"log"

	"path/filepath"

	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/viper"
)

// returns all pages of the hugo site located at path
func readSitePages(path string) hugolib.Pages {

	// warning: if you specify configFilename, it will not search path
	checkFatal(hugolib.LoadGlobalConfig(path, ""))

	// this would be done in InitializeConfig() that we're NOT calling
	// because it does not allow us to specify the source (path)
	dir, _ := filepath.Abs(path)
	viper.Set("WorkingDir", dir)

	sites, err := hugolib.NewHugoSitesFromConfiguration()

	if err != nil {
		log.Println("FATAL: Error creating sites", err)
	}

	if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		log.Println("FATAL: Error Processing Source Content", err)
	}

	return sites.Pages()
}

// checks if the page has a title (that will appear in the search result)
func pageHasTitle(page *hugolib.Page) (foundTitle bool) {
	foundTitle = len(page.Title) > 0
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
// 'Fails:section'		yes
// 'Folder1s:section'		yes
// 'Tag1:taxonomy'		no	dynamic content, do not index
// 'Tag2:taxonomy'		no	dynamic content, do not index
// 'Tags:taxonomyTerm'		no	dynamic content, do not index
// 'hugo-search test site:home'	yes
func pageHasValidContent(page *hugolib.Page) bool {
	switch page.Kind {
	case "page":
		if page.Title == "Search Results" {
			break
		}
		fallthrough
	case "section":
		fallthrough
	case "home":
		return true
	}
	if *verbose {
		log.Printf("Ignoring: %s [%s] of kind: %s", page.File.Path(), page.Title, page.Kind)
	}
	return false
}
