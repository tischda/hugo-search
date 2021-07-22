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

package babel

import (
	"bytes"
	"io"
	"path/filepath"
	"strconv"

	"github.com/cli/safeexec"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/resources/internal"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/pkg/errors"
)

// Options from https://babeljs.io/docs/en/options
type Options struct {
	Config string // Custom path to config file

	Minified   bool
	NoComments bool
	Compact    *bool
	Verbose    bool
	NoBabelrc  bool
}

func DecodeOptions(m map[string]interface{}) (opts Options, err error) {
	if m == nil {
		return
	}
	err = mapstructure.WeakDecode(m, &opts)
	return
}
func (opts Options) toArgs() []string {
	var args []string

	if opts.Minified {
		args = append(args, "--minified")
	}
	if opts.NoComments {
		args = append(args, "--no-comments")
	}
	if opts.Compact != nil {
		args = append(args, "--compact="+strconv.FormatBool(*opts.Compact))
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if opts.NoBabelrc {
		args = append(args, "--no-babelrc")
	}
	return args
}

// Client is the client used to do Babel transformations.
type Client struct {
	rs *resources.Spec
}

// New creates a new Client with the given specification.
func New(rs *resources.Spec) *Client {
	return &Client{rs: rs}
}

type babelTransformation struct {
	options Options
	rs      *resources.Spec
}

func (t *babelTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("babel", t.options)
}

// Transform shells out to babel-cli to do the heavy lifting.
// For this to work, you need some additional tools. To install them globally:
// npm install -g @babel/core @babel/cli
// If you want to use presets or plugins such as @babel/preset-env
// Then you should install those globally as well. e.g:
// npm install -g @babel/preset-env
// Instead of installing globally, you can also install everything as a dev-dependency (--save-dev instead of -g)
func (t *babelTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	const localBabelPath = "node_modules/.bin/"
	const binaryName = "babel"

	// Try first in the project's node_modules.
	csiBinPath := filepath.Join(t.rs.WorkingDir, localBabelPath, binaryName)

	binary := csiBinPath

	if _, err := safeexec.LookPath(binary); err != nil {
		// Try PATH
		binary = binaryName
		if _, err := safeexec.LookPath(binary); err != nil {

			// This may be on a CI server etc. Will fall back to pre-built assets.
			return herrors.ErrFeatureNotAvailable
		}
	}

	var configFile string
	logger := t.rs.Logger

	var errBuf bytes.Buffer
	infoW := loggers.LoggerToWriterWithPrefix(logger.Info(), "babel")

	if t.options.Config != "" {
		configFile = t.options.Config
	} else {
		configFile = "babel.config.js"
	}

	configFile = filepath.Clean(configFile)

	// We need an abolute filename to the config file.
	if !filepath.IsAbs(configFile) {
		configFile = t.rs.BaseFs.ResolveJSConfigFile(configFile)
		if configFile == "" && t.options.Config != "" {
			// Only fail if the user specificed config file is not found.
			return errors.Errorf("babel config %q not found:", configFile)
		}
	}

	var cmdArgs []string

	if configFile != "" {
		logger.Infoln("babel: use config file", configFile)
		cmdArgs = []string{"--config-file", configFile}
	}

	if optArgs := t.options.toArgs(); len(optArgs) > 0 {
		cmdArgs = append(cmdArgs, optArgs...)
	}
	cmdArgs = append(cmdArgs, "--filename="+ctx.SourcePath)

	cmd, err := hexec.SafeCommand(binary, cmdArgs...)
	if err != nil {
		return err
	}

	cmd.Stdout = ctx.To
	cmd.Stderr = io.MultiWriter(infoW, &errBuf)
	cmd.Env = hugo.GetExecEnviron(t.rs.WorkingDir, t.rs.Cfg, t.rs.BaseFs.Assets.Fs)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, ctx.From)
	}()

	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, errBuf.String())
	}

	return nil
}

// Process transforms the given Resource with the Babel processor.
func (c *Client) Process(res resources.ResourceTransformer, options Options) (resource.Resource, error) {
	return res.Transform(
		&babelTransformation{rs: c.rs, options: options},
	)
}
