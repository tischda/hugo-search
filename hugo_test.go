package main

import "testing"

const testHugoPath = "test"

// checks the number of pages built for the site (drafts don't count)
func TestReadSitePages(t *testing.T) {
	pages := readSitePages(testHugoPath)
	actual := pages.Len()
	expected := 5

	if expected != actual {
		var titles []string
		for _, page := range pages {
			titles = append(titles, "'"+page.Title()+":"+page.Kind()+"'")
		}
		t.Errorf("Expected: %d, was: %d, pages returned:\n%s", expected, actual, titles)
	}
}

// checks that pages with no title are correctly detected
func TestPageHasTitle(t *testing.T) {
	pages := readSitePages(testHugoPath)
	var a, b bool
	for _, page := range pages {
		if page.Title() == "Title-page-1" {
			a = pageHasTitle(page)
		} else if page.Title() == "" {
			b = !pageHasTitle(page)
		}
	}
	if !(a && b) {
		t.Errorf("Expected: has title==(true && false), was: (%v && %v)", a, b)
	}
}
