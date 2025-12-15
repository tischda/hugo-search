# Changelog

## [v1.4.0] - 15 December 2025

* Freshen-up build and release system
* Update hugo to 0.152.2 (fixes #1)
* Update bleve to 2.5.6 (had to take over bleve/http as well from v1)
* Use go modules instead of govendor
* Update javascript in test site

## [v1.3.0] - 2 December 2017

* Fix code to match latest versions of hugo and bleve

## [v1.2.0] - 11 January 2017

* Fix versions with govendor

## [v1.1.1] - 4 December 2016

* Code clean up
* Unregister index before closing
* Fix build due API changes in Hugo (page.RelPermalink())
* Do not index special pages such as taxonomies and search results
* Revert Makefile to work with make.exe from http://win-builds.org

## [v1.0.2] - 8 September 2016

* Fix build due to API changes

## [v1.0.1] - 9 August 2016

* Fix small bug in search.js
* Make searchUrl configurable

## [v1.0.0] - 3 April 2016

* First version
