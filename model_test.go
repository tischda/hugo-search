package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gohugoio/hugo/resources/page"
)

var expected *PageEntry

func setUp() {
	myDate, _ := time.Parse(time.RFC3339, "2015-12-09T22:15:11+01:00")
	expected = &PageEntry{
		Type:         "page",
		Section:      "",
		Content:      "Lorem ipsum Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n",
		WordCount:    10,
		ReadingTime:  1,
		Keywords:     []string{"keyword"},
		Date:         myDate,
		LastModified: myDate,
	}
}

// test single author
func TestNewPageForIndex(t *testing.T) {
	setUp()
	expected.Title = "Title-page-1"
	expected.Author = "Author1Page1"
	actual := newIndexEntry(findPage(expected.Title))
	comparePages(t, actual, expected)
}

// test multiple authors
func TestNewPageForIndex2(t *testing.T) {
	setUp()
	expected.Title = "Title-page-2"
	expected.Author = "Author1Page2, Author2Page2"
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
func findPage(title string) page.Page {
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
