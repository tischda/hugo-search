[![Build Status](https://github.com/tischda/hugo-search/actions/workflows/build.yml/badge.svg)](https://github.com/tischda/hugo-search/actions/workflows/build.yml)
[![Test Status](https://github.com/tischda/hugo-search/actions/workflows/test.yml/badge.svg)](https://github.com/tischda/hugo-search/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/tischda/hugo-search/badge.svg)](https://coveralls.io/r/tischda/hugo-search)
[![Linter Status](https://github.com/tischda/hugo-search/actions/workflows/linter.yml/badge.svg)](https://github.com/tischda/hugo-search/actions/workflows/linter.yml)
[![License](https://img.shields.io/github/license/tischda/hugo-search)](/LICENSE)
[![Release](https://img.shields.io/github/release/tischda/hugo-search.svg)](https://github.com/tischda/hugo-search/releases/latest)

# hugo search

A search engine for a [Hugo](http://gohugo.io) site using a [Bleve](http://www.blevesearch.com) index.

THIS REPOSITORY IS ARCHIVED.
I last tested hugo-search with Bleve v2.5.6 and Hugo v0.152.2 as a proof of concept (not for production use!).

## Install

~~~
go install https://github.com/tischda/hugo-search@latest
~~~

## Usage

~~~
Usage: hugo-search [OPTIONS]

OPTIONS:

  -a, --bindAddr"
          http listen address (default ":8080")
  -h, --hugoPath
          path to the hugo site (default ".")
  -i, --indexPath
          path to the bleve index (default "indexes/search.bleve")
  -vv, --verbose
          verbose output
  -?, --help
          display this help message
  -v, --version
          print version and exit
~~~


## Examples

Start hugo server:
~~~
hugo server --source test
~~~

In another console:
~~~
hugo-search.exe --hugoPath test
~~~

Open browser on [http://localhost:1313](http://localhost:1313)


## Query index

~~~
curl http://localhost:8080/api/search.bleve/_search -d "{\"query\":{\"query\":\"lorem\"}}"
~~~

result:
~~~
{"status":{"total":1,"failed":0,"successful":1},"hits":[{"index":"indexes/search.bleve","id":"/page1/","score":0.112627352112581,"sort":["_score"]},{"index":"indexes/search.bleve","id":"/page2/","score":0.11192996459618623,"sort":["_score"]},{"index":"indexes/search.bleve","id":"/parent1/page3/","score":0.10991304652804537,"sort":["_score"]}],"total_hits":3,"cost":2647,"max_score":0.112627352112581,"took":0,"facets":null}
~~~

## Explore index with bleve-explorer

Warning: Cannot use while `hugo-search` is running.

~~~
go install github.com/blevesearch/bleve-explorer@latest

bleve-explorer -dataDir indexes
~~~

Open browser on [http://localhost:8095/](http://localhost:8095/)
