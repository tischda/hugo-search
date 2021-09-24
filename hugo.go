package main

import (
	"log"

	"path/filepath"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/spf13/afero"
)

// returns all regular pages of the hugo site located at path
func readSitePages(path string) page.Pages {

	var sourceFs afero.Fs = hugofs.Os

	dir, err := filepath.Abs(path)
	exitOnError(err)

	config, _, err := hugolib.LoadConfig(hugolib.ConfigSourceDescriptor{
		Fs:         sourceFs,
		Path:       dir,
		WorkingDir: dir},
	)
	exitOnError(err)

	fs := hugofs.NewFrom(sourceFs, config)

	h, err := hugolib.NewHugoSites(deps.DepsCfg{Cfg: config, Fs: fs})
	exitOnError(err)

	err = h.Build(hugolib.BuildCfg{SkipRender: true})
	exitOnError(err)

	// current language only, for all languages use AllRegularPages()
	return h.Sites[0].RegularPages()

	// TODO: does not include the static home page, for this we could use AllPages but this is too much
}

// checks if the page has a title, which is required to be displayed in the search result
func pageHasTitle(p page.Page) (foundTitle bool) {
	foundTitle = len(p.Title()) > 0
	if !foundTitle && *verbose {
		log.Println("WARN: Title is missing in file metadata:", p.File().Path())
	}
	return
}
