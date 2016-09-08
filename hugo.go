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

// checks if a page has a title (which will appear in the search result)
func pageHasTitle(page *hugolib.Page) (foundTitle bool) {
	foundTitle = len(page.Title) > 0
	if !foundTitle && *verbose {
		log.Println("WARN: Title metadata missing in document:", page.File.Path())
	}
	return
}
