// Copyright 2018 The Hugo Authors. All rights reserved.
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

package hugofs

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

var (
	_ afero.Fs      = (*languageCompositeFs)(nil)
	_ afero.Lstater = (*languageCompositeFs)(nil)
)

type languageCompositeFs struct {
	*afero.CopyOnWriteFs
}

// NewLanguageCompositeFs creates a composite and language aware filesystem.
// This is a hybrid filesystem. To get a specific file in Open, Stat etc., use the full filename
// to the target filesystem. This information is available in Readdir, Stat etc. via the
// special LanguageFileInfo FileInfo implementation.
func NewLanguageCompositeFs(base, overlay afero.Fs) afero.Fs {
	return &languageCompositeFs{afero.NewCopyOnWriteFs(base, overlay).(*afero.CopyOnWriteFs)}
}

// Open takes the full path to the file in the target filesystem. If it is a directory, it gets merged
// using the language as a weight.
func (fs *languageCompositeFs) Open(name string) (afero.File, error) {
	f, err := fs.CopyOnWriteFs.Open(name)
	if err != nil {
		return nil, err
	}

	fu, ok := f.(*afero.UnionFile)
	if ok {
		// This is a directory: Merge it.
		fu.Merger = LanguageDirsMerger
	}
	return f, nil
}

// LanguageDirsMerger implements the afero.DirsMerger interface, which is used
// to merge two directories.
var LanguageDirsMerger = func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error) {
	m := make(map[string]FileMetaInfo)

	getKey := func(fim FileMetaInfo) string {
		return path.Join(fim.Meta().Lang(), fim.Name())
	}

	for _, fi := range lofi {
		fim := fi.(FileMetaInfo)
		m[getKey(fim)] = fim
	}

	for _, fi := range bofi {
		fim := fi.(FileMetaInfo)
		key := getKey(fim)
		_, found := m[key]
		if !found {
			m[key] = fim
		}
	}

	merged := make([]os.FileInfo, len(m))
	i := 0
	for _, v := range m {
		merged[i] = v
		i++
	}

	return merged, nil
}
