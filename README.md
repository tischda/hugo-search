# hugo-search [![Build Status](https://travis-ci.org/tischda/hugo-search.svg?branch=master)](https://travis-ci.org/tischda/hugo-search)

A [Bleve](http://www.blevesearch.com) search server for your [Hugo](http://gohugo.io) site.

### Install

~~~
go get github.com/tischda/hugo-search
~~~

### Usage

~~~
Usage of hugo-search:
  -addr string
        http listen address (default ":8080")
  -hugoPath string
        path of the hugo site (default ".")
  -indexPath string
        path of the bleve index (default "indexes/search.bleve")
  -stepAnalysis
        display memory and timing of different steps of the program
  -verbose    verbose output
  -version
        print version and exit
~~~

### Query index

~~~
$ curl http://localhost:8080/api/search.bleve/_search -d '{"query":{"query":"lorem"}}'
{"request":{"query":{"query":"lorem","boost":1},"size":0,"from":0,"highlight":null,"fields":null,"facets":null,"explain":false},"hits":[],"total_hits":2,"max_score":0.15713484143442302,"took":0,"facets":{}}

{"status":{"total":1,"failed":0,"successful":1},"request":{"query":{"query":"lorem","boost":1},"size":0,"from":0,"highlight":null,"fields":null,"facets":null,"explain":false},"hits":[],"total_hits":3,"max_score":0.15713484143442302,"took":0,"facets":{}}
~~~

### bleve-explorer

Warning: Cannot use while `hugo-search` is running.

~~~
go get github.com/blevesearch/bleve-explorer

start bleve-explorer -dataDir indexes
start http://localhost:8095/
~~~
