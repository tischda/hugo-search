module github.com/tischda/hugo-search

go 1.25

require (
	github.com/BurntSushi/locker v0.0.0-20171006230638-a6e239ea1c69
	github.com/JohannesKaufmann/dom v0.2.0
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.5.0
	github.com/RoaringBitmap/roaring v1.9.4
	github.com/alecthomas/chroma/v2 v2.20.0
	github.com/armon/go-radix v1.0.1-0.20221118154546-54df44f2176c
	github.com/aymerick/douceur v0.2.0
	github.com/bep/clocks v0.5.0
	github.com/bep/debounce v1.2.1
	github.com/bep/gitmap v1.9.0
	github.com/bep/goat v0.5.0
	github.com/bep/godartsass/v2 v2.5.0
	github.com/bep/golibsass v1.2.0
	github.com/bep/goportabletext v0.1.0
	github.com/bep/gowebp v0.4.0
	github.com/bep/helpers v0.6.0
	github.com/bep/imagemeta v0.12.0
	github.com/bep/lazycache v0.8.0
	github.com/bep/logg v0.4.0
	github.com/bep/overlayfs v0.10.0
	github.com/bep/tmc v0.5.1
	github.com/bits-and-blooms/bitset v1.24.4
	github.com/blevesearch/bleve v1.0.14
	github.com/blevesearch/bleve/v2 v2.5.6
	github.com/blevesearch/go-porterstemmer v1.0.3
	github.com/blevesearch/mmap-go v1.0.4
	github.com/blevesearch/segment v0.9.1
	github.com/blevesearch/snowballstem v0.9.0
	github.com/blevesearch/zap/v11 v11.0.14
	github.com/blevesearch/zap/v12 v12.0.14
	github.com/blevesearch/zap/v13 v13.0.6
	github.com/blevesearch/zap/v14 v14.0.5
	github.com/blevesearch/zap/v15 v15.0.3
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/clbanning/mxj/v2 v2.7.0
	github.com/clipperhouse/displaywidth v0.6.2
	github.com/clipperhouse/stringish v0.1.1
	github.com/clipperhouse/uax29/v2 v2.3.0
	github.com/couchbase/vellum v1.0.2
	github.com/disintegration/gift v1.2.1
	github.com/dlclark/regexp2 v1.11.5
	github.com/evanw/esbuild v0.27.1
	github.com/fatih/color v1.18.0
	github.com/frankban/quicktest v1.14.6
	github.com/fsnotify/fsnotify v1.9.0
	github.com/getkin/kin-openapi v0.133.0
	github.com/go-openapi/jsonpointer v0.22.4
	github.com/go-openapi/swag/jsonname v0.25.4
	github.com/gobuffalo/flect v1.0.3
	github.com/gobwas/glob v0.2.3
	github.com/goccy/go-yaml v1.19.0
	github.com/gohugoio/go-i18n/v2 v2.1.3-0.20251018145728-cfcc22d823c6
	github.com/gohugoio/hashstructure v0.6.0
	github.com/gohugoio/httpcache v0.8.0
	github.com/gohugoio/hugo v0.152.2
	github.com/gohugoio/hugo-goldmark-extensions/extras v0.5.0
	github.com/gohugoio/hugo-goldmark-extensions/passthrough v0.3.1
	github.com/gohugoio/locales v0.14.0
	github.com/gohugoio/localescompressed v1.0.1
	github.com/golang/protobuf v1.5.4
	github.com/golang/snappy v1.0.0
	github.com/google/go-cmp v0.7.0
	github.com/gorilla/css v1.0.1
	github.com/hairyhenderson/go-codeowners v0.7.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/jdkato/prose v1.2.1
	github.com/josharian/intern v1.0.0
	github.com/kr/pretty v0.3.1
	github.com/kr/text v0.2.0
	github.com/kyokomi/emoji/v2 v2.2.13
	github.com/mailru/easyjson v0.9.1
	github.com/makeworld-the-better-one/dither/v2 v2.4.0
	github.com/marekm4/color-extractor v1.2.1
	github.com/mattn/go-colorable v0.1.14
	github.com/mattn/go-isatty v0.0.20
	github.com/mattn/go-runewidth v0.0.19
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/mschoch/smat v0.2.0
	github.com/muesli/smartcrop v0.3.0
	github.com/niklasfasching/go-org v1.9.1
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6
	github.com/olekukonko/errors v1.1.0
	github.com/olekukonko/ll v0.1.3
	github.com/olekukonko/tablewriter v1.1.2
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/perimeterx/marshmallow v1.1.5
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.14.1
	github.com/rs/cors v1.11.1
	github.com/spf13/afero v1.15.0
	github.com/spf13/cast v1.10.0
	github.com/steveyen/gtreap v0.1.0
	github.com/tdewolff/minify/v2 v2.24.8
	github.com/tdewolff/parse/v2 v2.8.5
	github.com/tetratelabs/wazero v1.10.1
	github.com/willf/bitset v1.1.11
	github.com/woodsbury/decimal128 v1.4.0
	github.com/yuin/goldmark v1.7.13
	github.com/yuin/goldmark-emoji v1.0.6
	go.etcd.io/bbolt v1.4.3
	golang.org/x/image v0.34.0
	golang.org/x/mod v0.31.0
	golang.org/x/net v0.48.0
	golang.org/x/sync v0.19.0
	golang.org/x/sys v0.39.0
	golang.org/x/text v0.32.0
	golang.org/x/tools v0.40.0
	google.golang.org/protobuf v1.36.11
	rsc.io/qr v0.2.0
)

require (
	github.com/RoaringBitmap/roaring/v2 v2.14.4 // indirect
	github.com/blevesearch/bleve_index_api v1.2.11 // indirect
	github.com/blevesearch/geo v0.2.4 // indirect
	github.com/blevesearch/go-faiss v1.0.26 // indirect
	github.com/blevesearch/gtreap v0.1.1 // indirect
	github.com/blevesearch/scorch_segment_api/v2 v2.3.13 // indirect
	github.com/blevesearch/upsidedown_store_api v1.0.2 // indirect
	github.com/blevesearch/vellum v1.1.0 // indirect
	github.com/blevesearch/zapx/v11 v11.4.2 // indirect
	github.com/blevesearch/zapx/v12 v12.4.2 // indirect
	github.com/blevesearch/zapx/v13 v13.4.2 // indirect
	github.com/blevesearch/zapx/v14 v14.4.2 // indirect
	github.com/blevesearch/zapx/v15 v15.4.2 // indirect
	github.com/blevesearch/zapx/v16 v16.2.8 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
)
