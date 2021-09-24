package main

import (
	"strings"
	"time"

	"github.com/gohugoio/hugo/resources/page"
)

// PageEntry maps the hugo internal page structure to a JSON structure
// that blevesearch can understand.
type PageEntry struct {
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	Section      string    `json:"section"`
	Content      string    `json:"content"`
	WordCount    float64   `json:"word_count"`
	ReadingTime  float64   `json:"reading_time"`
	Keywords     []string  `json:"keywords"`
	Date         time.Time `json:"date"`
	LastModified time.Time `json:"last_modified"`
	Author       string    `json:"author"`
}

func newIndexEntry(p page.Page) *PageEntry {
	var author string

	switch str := p.Params()["author"].(type) {
	case string:
		author = str
	case []string:
		author = strings.Join(str, ", ")
	}

	return &PageEntry{
		Title:        p.Title(),
		Type:         p.Type(),
		Section:      p.Section(),
		Content:      p.Plain(),
		WordCount:    float64(p.WordCount()),
		ReadingTime:  float64(p.ReadingTime()),
		Keywords:     p.Keywords(),
		Date:         p.Date(),
		LastModified: p.Lastmod(),
		Author:       author,
	}
}
