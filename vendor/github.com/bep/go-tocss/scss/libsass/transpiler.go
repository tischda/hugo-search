// Copyright © 2018 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package libsass a SCSS transpiler to CSS using github.com/wellington/go-libsass/libs.
package libsass

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bep/go-tocss/scss"
	"github.com/bep/go-tocss/tocss"
	"github.com/wellington/go-libsass/libs"
)

type libsassTranspiler struct {
	options scss.Options
}

// New creates a new libsass transpiler configured with the given options.
func New(options scss.Options) (tocss.Transpiler, error) {
	return &libsassTranspiler{options: options}, nil
}

// Execute transpiles the SCSS from src into dst. Note that you can import
// older SASS (.sass) files, but the main entry (src) currently needs to be SCSS.
func (t *libsassTranspiler) Execute(dst io.Writer, src io.Reader) (tocss.Result, error) {
	var result tocss.Result
	var sourceStr string

	if t.options.SassSyntax {
		// LibSass does not support this directly, so have to handle the main SASS content
		// special.
		var buf bytes.Buffer
		err := libs.ToScss(src, &buf)
		if err != nil {
			return result, err
		}
		sourceStr = buf.String()
	} else {
		b, err := ioutil.ReadAll(src)
		if err != nil {
			return result, err
		}
		sourceStr = string(b)
	}

	dataCtx := libs.SassMakeDataContext(sourceStr)

	opts := libs.SassDataContextGetOptions(dataCtx)

	{
		// Set options

		if t.options.ImportResolver != nil {
			idx := libs.BindImporter(opts, t.options.ImportResolver)
			defer libs.RemoveImporter(idx)
		}

		if t.options.Precision != 0 {
			libs.SassOptionSetPrecision(opts, t.options.Precision)
		}

		if t.options.SourceMapFilename != "" {
			libs.SassOptionSetSourceMapFile(opts, t.options.SourceMapFilename)
		}

		if t.options.SourceMapRoot != "" {
			libs.SassOptionSetSourceMapRoot(opts, t.options.SourceMapRoot)
		}

		if t.options.OutputPath != "" {
			libs.SassOptionSetOutputPath(opts, t.options.OutputPath)
		}
		if t.options.InputPath != "" {
			libs.SassOptionSetInputPath(opts, t.options.InputPath)
		}

		libs.SassOptionSetSourceMapContents(opts, t.options.SourceMapContents)
		libs.SassOptionSetOmitSourceMapURL(opts, t.options.OmitSourceMapURL)
		libs.SassOptionSetSourceMapEmbed(opts, t.options.EnableEmbeddedSourceMap)
		libs.SassOptionSetIncludePath(opts, strings.Join(t.options.IncludePaths, string(os.PathListSeparator)))
		libs.SassOptionSetOutputStyle(opts, int(t.options.OutputStyle))
		libs.SassOptionSetSourceComments(opts, false)
		libs.SassDataContextSetOptions(dataCtx, opts)
	}

	ctx := libs.SassDataContextGetContext(dataCtx)
	compiler := libs.SassMakeDataCompiler(dataCtx)

	libs.SassCompilerParse(compiler)
	libs.SassCompilerExecute(compiler)

	defer libs.SassDeleteCompiler(compiler)

	outputString := libs.SassContextGetOutputString(ctx)

	io.WriteString(dst, outputString)

	if status := libs.SassContextGetErrorStatus(ctx); status != 0 {
		return result, scss.JSONToError(libs.SassContextGetErrorJSON(ctx))
	}

	result.SourceMapFilename = libs.SassOptionGetSourceMapFile(opts)
	result.SourceMapContent = libs.SassContextGetSourceMapString(ctx)

	return result, nil
}
