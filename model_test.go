package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gohugoio/hugo/hugolib"
)

var expected *PageEntry

func setUp() {
	myDate, _ := time.Parse(time.RFC3339, "2015-12-09T22:15:11+01:00")
	expected = &PageEntry{
		Type:         "page",
		Section:      "",
		Content:      " Lorem ipsum Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n",
		WordCount:    10,
		ReadingTime:  1,
		Keywords:     []string{},
		Date:         myDate,
		LastModified: myDate,
	}
}

// test single author
func TestNewPageForIndex(t *testing.T) {
	setUp()
	expected.Title = "title-page-1"
	expected.Author = "author-page-1"
	actual := newIndexEntry(findPage(expected.Title))
	comparePages(t, actual, expected)
}

// test multiple authors
func TestNewPageForIndex2(t *testing.T) {
	setUp()
	expected.Title = "title-page-2"
	expected.Author = "author-1-page-2, author-2-page-2"
	actual := newIndexEntry(findPage(expected.Title))
	comparePages(t, actual, expected)
}

func comparePages(t *testing.T, actual *PageEntry, expected *PageEntry) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Values don't match!")
		printPage("Expected:", expected)
		printPage("but was:", actual)
	}
}

// find first page with specified title
func findPage(title string) *hugolib.Page {
	pages := readSitePages(testHugoPath)
	for _, page := range pages {
		if page.Title() == title {
			return page
		}
	}
	return nil
}

// print struct as pretty formatted json
func printPage(label string, page *PageEntry) {
	fmt.Println(label)
	data, _ := json.MarshalIndent(page, "", "    ")
	os.Stdout.Write(data)
	fmt.Println()
}
