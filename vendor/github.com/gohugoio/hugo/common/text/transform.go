// Copyright 2019 The Hugo Authors. All rights reserved.
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

package text

import (
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var accentTransformerPool = &sync.Pool{
	New: func() interface{} {
		return transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	},
}

// RemoveAccents removes all accents from b.
func RemoveAccents(b []byte) []byte {
	t := accentTransformerPool.Get().(transform.Transformer)
	b, _, _ = transform.Bytes(t, b)
	t.Reset()
	accentTransformerPool.Put(t)
	return b
}

// RemoveAccentsString removes all accents from s.
func RemoveAccentsString(s string) string {
	t := accentTransformerPool.Get().(transform.Transformer)
	s, _, _ = transform.String(t, s)
	t.Reset()
	accentTransformerPool.Put(t)
	return s
}
