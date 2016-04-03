package main

import (
	"log"

	"path/filepath"

	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/viper"
)

// returns all pages of the hugo site located at path
func readSitePages(path string) hugolib.Pages {
	initializeConfig(path)
	site := &hugolib.Site{}
	checkFatal(site.Process())
	return site.Pages
}

// Using hugolib.InitializeConfig() is cumbersome. In this stripped down
// version, all we do here is setting a different source path.
func initializeConfig(path string) {
	dir, _ := filepath.Abs(path)
	viper.Set("WorkingDir", dir)

	// TODO: in the original code, but is this really needed here ?
	// load the configuration file from disk and register aliases
	viper.AddConfigPath(path)
	checkFatal(viper.ReadInConfig())
	viper.RegisterAlias("indexes", "taxonomies")

	// could we set this in the config file instead ?
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("BuildDrafts", false)
	viper.SetDefault("BuildFuture", false)
}

// checks if a page has a title (which will appear in the search result)
func pageHasTitle(page *hugolib.Page) (foundTitle bool) {
	foundTitle = len(page.Title) > 0
	if !foundTitle && *verbose {
		log.Println("WARN: Title metadata missing in document:", page.File.Path())
	}
	return
}
