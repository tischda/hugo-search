package main

import "testing"

const testHugoPath = "test"

// checks number of posts in site (drafts don't count)
func TestReadSitePages(t *testing.T) {
	actual := readSitePages(testHugoPath).Len()
	expected := 5

	if expected != actual {
		t.Errorf("Expected: %d, was: %d", expected, actual)
	}
}

// checks that pages with no title are correctly detected
func TestPageHasTitle(t *testing.T) {
	pages := readSitePages(testHugoPath)
	var a, b bool
	for _, page := range pages {
		if page.Title == "title-page-1" {
			a = pageHasTitle(page)
		} else if page.Title == "" {
			b = !pageHasTitle(page)
		}
	}
	if !(a && b) {
		t.Errorf("Expected: has title==(true && false), was: (%v && %v)", a, b)
	}
}
