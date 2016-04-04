package main

import (
	"log"

	"path/filepath"

	"github.com/spf13/hugo/helpers"
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

// stripped down version of hugolib.InitializeConfig()
func initializeConfig(path string) {
	dir, _ := filepath.Abs(path)
	viper.Set("WorkingDir", dir)

	viper.AddConfigPath(path)
	checkFatal(viper.ReadInConfig())
	viper.RegisterAlias("indexes", "taxonomies")

	loadDefaultSettings()
}

// copy of commands.loadDefaultSettings() since not public anymore
func loadDefaultSettings() {
	viper.SetDefault("cleanDestinationDir", false)
	viper.SetDefault("Watch", false)
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("DisableRSS", false)
	viper.SetDefault("DisableSitemap", false)
	viper.SetDefault("DisableRobotsTXT", false)
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("BuildDrafts", false)
	viper.SetDefault("BuildFuture", false)
	viper.SetDefault("UglyURLs", false)
	viper.SetDefault("Verbose", false)
	viper.SetDefault("IgnoreCache", false)
	viper.SetDefault("CanonifyURLs", false)
	viper.SetDefault("RelativeURLs", false)
	viper.SetDefault("RemovePathAccents", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	viper.SetDefault("Permalinks", make(hugolib.PermalinkOverrides, 0))
	viper.SetDefault("Sitemap", hugolib.Sitemap{Priority: -1, Filename: "sitemap.xml"})
	viper.SetDefault("DefaultExtension", "html")
	viper.SetDefault("PygmentsStyle", "monokai")
	viper.SetDefault("PygmentsUseClasses", false)
	viper.SetDefault("PygmentsCodeFences", false)
	viper.SetDefault("PygmentsOptions", "")
	viper.SetDefault("DisableLiveReload", false)
	viper.SetDefault("PluralizeListTitles", true)
	viper.SetDefault("PreserveTaxonomyNames", false)
	viper.SetDefault("ForceSyncStatic", false)
	viper.SetDefault("FootnoteAnchorPrefix", "")
	viper.SetDefault("FootnoteReturnLinkContents", "")
	viper.SetDefault("NewContentEditor", "")
	viper.SetDefault("Paginate", 10)
	viper.SetDefault("PaginatePath", "page")
	viper.SetDefault("Blackfriday", helpers.NewBlackfriday())
	viper.SetDefault("RSSUri", "index.xml")
	viper.SetDefault("SectionPagesMenu", "")
	viper.SetDefault("DisablePathToLower", false)
	viper.SetDefault("HasCJKLanguage", false)
	viper.SetDefault("EnableEmoji", false)
}

// checks if a page has a title (which will appear in the search result)
func pageHasTitle(page *hugolib.Page) (foundTitle bool) {
	foundTitle = len(page.Title) > 0
	if !foundTitle && *verbose {
		log.Println("WARN: Title metadata missing in document:", page.File.Path())
	}
	return
}
