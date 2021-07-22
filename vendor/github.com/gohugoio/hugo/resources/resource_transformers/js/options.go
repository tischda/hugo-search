// Copyright 2020 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package js

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/media"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

const (
	nsImportHugo = "ns-hugo"
	nsParams     = "ns-params"

	stdinImporter = "<stdin>"
)

// Options esbuild configuration
type Options struct {
	// If not set, the source path will be used as the base target path.
	// Note that the target path's extension may change if the target MIME type
	// is different, e.g. when the source is TypeScript.
	TargetPath string

	// Whether to minify to output.
	Minify bool

	// Whether to write mapfiles
	SourceMap string

	// The language target.
	// One of: es2015, es2016, es2017, es2018, es2019, es2020 or esnext.
	// Default is esnext.
	Target string

	// The output format.
	// One of: iife, cjs, esm
	// Default is to esm.
	Format string

	// External dependencies, e.g. "react".
	Externals []string `hash:"set"`

	// User defined symbols.
	Defines map[string]interface{}

	// User defined params. Will be marshaled to JSON and available as "@params", e.g.
	//     import * as params from '@params';
	Params interface{}

	// What to use instead of React.createElement.
	JSXFactory string

	// What to use instead of React.Fragment.
	JSXFragment string

	// There is/was a bug in WebKit with severe performance issue with the tracking
	// of TDZ checks in JavaScriptCore.
	//
	// Enabling this flag removes the TDZ and `const` assignment checks and
	// may improve performance of larger JS codebases until the WebKit fix
	// is in widespread use.
	//
	// See https://bugs.webkit.org/show_bug.cgi?id=199866
	// Deprecated: This no longer have any effect and will be removed.
	// TODO(bep) remove. See https://github.com/evanw/esbuild/commit/869e8117b499ca1dbfc5b3021938a53ffe934dba
	AvoidTDZ bool

	mediaType  media.Type
	outDir     string
	contents   string
	sourcefile string
	resolveDir string
	tsConfig   string
}

func decodeOptions(m map[string]interface{}) (Options, error) {
	var opts Options

	if err := mapstructure.WeakDecode(m, &opts); err != nil {
		return opts, err
	}

	if opts.TargetPath != "" {
		opts.TargetPath = helpers.ToSlashTrimLeading(opts.TargetPath)
	}

	opts.Target = strings.ToLower(opts.Target)
	opts.Format = strings.ToLower(opts.Format)

	return opts, nil
}

var extensionToLoaderMap = map[string]api.Loader{
	".js":   api.LoaderJS,
	".mjs":  api.LoaderJS,
	".cjs":  api.LoaderJS,
	".jsx":  api.LoaderJSX,
	".ts":   api.LoaderTS,
	".tsx":  api.LoaderTSX,
	".css":  api.LoaderCSS,
	".json": api.LoaderJSON,
	".txt":  api.LoaderText,
}

func loaderFromFilename(filename string) api.Loader {
	l, found := extensionToLoaderMap[filepath.Ext(filename)]
	if found {
		return l
	}
	return api.LoaderJS
}

func createBuildPlugins(c *Client, opts Options) ([]api.Plugin, error) {
	fs := c.rs.Assets

	resolveImport := func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		isStdin := args.Importer == stdinImporter
		var relDir string
		if !isStdin {
			rel, found := fs.MakePathRelative(args.Importer)
			if !found {
				// Not in any of the /assets folders.
				// This is an import from a node_modules, let
				// ESBuild resolve this.
				return api.OnResolveResult{}, nil
			}
			relDir = filepath.Dir(rel)
		} else {
			relDir = filepath.Dir(opts.sourcefile)
		}

		impPath := args.Path

		// Imports not starting with a "." is assumed to live relative to /assets.
		// Hugo makes no assumptions about the directory structure below /assets.
		if relDir != "" && strings.HasPrefix(impPath, ".") {
			impPath = filepath.Join(relDir, impPath)
		}

		findFirst := func(base string) hugofs.FileMeta {
			// This is the most common sub-set of ESBuild's default extensions.
			// We assume that imports of JSON, CSS etc. will be using their full
			// name with extension.
			for _, ext := range []string{".js", ".ts", ".tsx", ".jsx"} {
				if fi, err := fs.Fs.Stat(base + ext); err == nil {
					return fi.(hugofs.FileMetaInfo).Meta()
				}
			}

			// Not found.
			return nil
		}

		var m hugofs.FileMeta

		// First the path as is.
		fi, err := fs.Fs.Stat(impPath)

		if err == nil {
			if fi.IsDir() {
				m = findFirst(filepath.Join(impPath, "index"))
			} else {
				m = fi.(hugofs.FileMetaInfo).Meta()
			}
		} else {
			// It may be a regular file imported without an extension.
			m = findFirst(impPath)
		}
		//

		if m != nil {
			// Store the source root so we can create a jsconfig.json
			// to help intellisense when the build is done.
			// This should be a small number of elements, and when
			// in server mode, we may get stale entries on renames etc.,
			// but that shouldn't matter too much.
			c.rs.JSConfigBuilder.AddSourceRoot(m.SourceRoot())
			return api.OnResolveResult{Path: m.Filename(), Namespace: nsImportHugo}, nil
		}

		// Fall back to ESBuild's resolve.
		return api.OnResolveResult{}, nil
	}

	importResolver := api.Plugin{
		Name: "hugo-import-resolver",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `.*`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return resolveImport(args)

				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: nsImportHugo},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					b, err := ioutil.ReadFile(args.Path)

					if err != nil {
						return api.OnLoadResult{}, errors.Wrapf(err, "failed to read %q", args.Path)
					}
					c := string(b)
					return api.OnLoadResult{
						// See https://github.com/evanw/esbuild/issues/502
						// This allows all modules to resolve dependencies
						// in the main project's node_modules.
						ResolveDir: opts.resolveDir,
						Contents:   &c,
						Loader:     loaderFromFilename(args.Path),
					}, nil
				})
		},
	}

	params := opts.Params
	if params == nil {
		// This way @params will always resolve to something.
		params = make(map[string]interface{})
	}

	b, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal params")
	}
	bs := string(b)
	paramsPlugin := api.Plugin{
		Name: "hugo-params-plugin",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: `^@params$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return api.OnResolveResult{
						Path:      args.Path,
						Namespace: nsParams,
					}, nil
				})
			build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: nsParams},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					return api.OnLoadResult{
						Contents: &bs,
						Loader:   api.LoaderJSON,
					}, nil
				})
		},
	}

	return []api.Plugin{importResolver, paramsPlugin}, nil

}

func toBuildOptions(opts Options) (buildOptions api.BuildOptions, err error) {

	var target api.Target
	switch opts.Target {
	case "", "esnext":
		target = api.ESNext
	case "es5":
		target = api.ES5
	case "es6", "es2015":
		target = api.ES2015
	case "es2016":
		target = api.ES2016
	case "es2017":
		target = api.ES2017
	case "es2018":
		target = api.ES2018
	case "es2019":
		target = api.ES2019
	case "es2020":
		target = api.ES2020
	default:
		err = fmt.Errorf("invalid target: %q", opts.Target)
		return
	}

	mediaType := opts.mediaType
	if mediaType.IsZero() {
		mediaType = media.JavascriptType
	}

	var loader api.Loader
	switch mediaType.SubType {
	// TODO(bep) ESBuild support a set of other loaders, but I currently fail
	// to see the relevance. That may change as we start using this.
	case media.JavascriptType.SubType:
		loader = api.LoaderJS
	case media.TypeScriptType.SubType:
		loader = api.LoaderTS
	case media.TSXType.SubType:
		loader = api.LoaderTSX
	case media.JSXType.SubType:
		loader = api.LoaderJSX
	default:
		err = fmt.Errorf("unsupported Media Type: %q", opts.mediaType)
		return
	}

	var format api.Format
	// One of: iife, cjs, esm
	switch opts.Format {
	case "", "iife":
		format = api.FormatIIFE
	case "esm":
		format = api.FormatESModule
	case "cjs":
		format = api.FormatCommonJS
	default:
		err = fmt.Errorf("unsupported script output format: %q", opts.Format)
		return
	}

	var defines map[string]string
	if opts.Defines != nil {
		defines = cast.ToStringMapString(opts.Defines)
	}

	// By default we only need to specify outDir and no outFile
	var outDir = opts.outDir
	var outFile = ""
	var sourceMap api.SourceMap
	switch opts.SourceMap {
	case "inline":
		sourceMap = api.SourceMapInline
	case "":
		sourceMap = api.SourceMapNone
	default:
		err = fmt.Errorf("unsupported sourcemap type: %q", opts.SourceMap)
		return
	}

	buildOptions = api.BuildOptions{
		Outfile: outFile,
		Bundle:  true,

		Target:    target,
		Format:    format,
		Sourcemap: sourceMap,

		MinifyWhitespace:  opts.Minify,
		MinifyIdentifiers: opts.Minify,
		MinifySyntax:      opts.Minify,

		Outdir: outDir,
		Define: defines,

		External: opts.Externals,

		JSXFactory:  opts.JSXFactory,
		JSXFragment: opts.JSXFragment,

		AvoidTDZ: opts.AvoidTDZ,

		Tsconfig: opts.tsConfig,

		Stdin: &api.StdinOptions{
			Contents:   opts.contents,
			Sourcefile: opts.sourcefile,
			ResolveDir: opts.resolveDir,
			Loader:     loader,
		},
	}
	return

}
