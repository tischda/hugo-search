package main

import (
	"fmt"
	"path/filepath"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/allconfig"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/resources/page"
)

func readSitePages(path string) page.Pages {

	sourceFs := hugofs.Os

	dir, err := filepath.Abs(path)
	exitOnError(err)

	params := make(map[string]interface{})

	var flags config.Provider = config.NewFrom(params)
	flags.Set("workingDir", dir)

	configs, err := allconfig.LoadConfig(
		allconfig.ConfigSourceDescriptor{
			Flags:       flags,
			Fs:          sourceFs,
			Filename:    "config.toml",
			ConfigDir:   ".",
			Environment: "production",
		},
	)
	exitOnError(err)

	conf := configs.GetFirstLanguageConfig().BaseConfig()
	fs := hugofs.NewFrom(sourceFs, conf)
	sites, err := hugolib.NewHugoSites(deps.DepsCfg{
		Fs:      fs,
		Configs: configs,
	})
	exitOnError(err)

	err = sites.Build(hugolib.BuildCfg{SkipRender: true})
	exitOnError(err)

	fmt.Printf("All pages: %d, regular pages: %d\n", len(sites.Pages()), len(sites.Sites[0].RegularPages()))

	return sites.Sites[0].RegularPages()
}
