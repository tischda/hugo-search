// Copyright 2025 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package goportabletext

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/bep/goportabletext/internal/portabletext"
)

const (
	newline   = "\n"
	indentOne = "   "
)

var newlineb = []byte(newline)

// ToMarkdown converts the given Portable Text blocks in opts.Src to Markdown and writes the result to opts.Dst.
func ToMarkdown(opts ToMarkdownOptions) error {
	blocks, ok := opts.Src.(portabletext.Blocks)
	if !ok {
		var err error
		blocks, err = portabletext.Parse(opts.Src)
		if err != nil {
			return err
		}
	}
	mw := getWriter()
	mw.dst = opts.Dst
	mw.src = blocks
	defer putWriter(mw)

	return mw.write()
}

// ToMarkdownOptions provides options for the ToMarkdown function.
type ToMarkdownOptions struct {
	// The destination writer.
	Dst io.Writer

	// The source Portable Text blocks.
	// Can be either
	// * a io.Reader with a JSON value or array of blocks.
	// * map[string]any or a []any.
	// * a portabletext.Blocks slice (used in benchmarks only).
	Src any
}

type markdownWriter struct {
	dst io.Writer
	src portabletext.Blocks
	err error

	// parser state.
	currentListItem string
	markDefs        map[string]portabletext.MarkDef
	marksOpen       [][]string
	marksClose      [][]string
}

func (m *markdownWriter) setListState(block portabletext.Block) {
	if block.ListItem == "" {
		m.currentListItem = ""
	}

	if block.Level == 1 {
		if block.ListItem != "" && (m.currentListItem == "" || m.currentListItem != block.ListItem) {
			m.currentListItem = block.ListItem
		}
	}
}

func (m *markdownWriter) indent(level int) string {
	if level <= 1 {
		return ""
	}
	return strings.Repeat(indentOne, level-1)
}

func (m *markdownWriter) indentText(s string, level int) string {
	if level == 1 {
		return s
	}
	indent := m.indent(level)
	if m.currentListItem != "" {
		// Get the indentation in line with the list item.
		switch m.currentListItem {
		case "number":
			indent += "   "
		default:
			indent += "  "
		}
	}
	i := strings.Index(s, "\n")
	if i == -1 {
		return s
	}

	var k int
	var sb strings.Builder
	for {

		j := strings.IndexFunc(s[i+k+1:], func(r rune) bool {
			return r != '\n'
		})

		if j == -1 {
			break
		}

		sb.WriteString(s[k : k+i+j+1])
		sb.WriteString(indent)

		k += i + j + 1

		i = strings.Index(s[k:], "\n")
		if i == -1 {
			break
		}

	}

	if k < len(s) {
		sb.WriteString(s[k:])
	}

	return sb.String()
}

// We always insert a newline after a block. This method writes a newline
// before a block if needed.
func (m *markdownWriter) newlineBeforeBlockIfNeeded(i int) {
	if i == 0 {
		return
	}
	prev := m.src[i-1]
	current := m.src[i]

	if func() bool {
		if !current.HasText() {
			return false
		}
		// E.g. from paragraph to list.
		if prev.ListItem == "" && current.ListItem != "" {
			return true
		}
		// New list.
		if current.Level == 1 && current.ListItem != "" && m.currentListItem != current.ListItem {
			return true
		}
		if prev.ListItem != "" && current.ListItem != "" {
			return false
		}

		return true
	}() {
		m.writeNewline()
	}
}

func (m *markdownWriter) write() error {
	if len(m.src) == 0 {
		m.writeNewline()
		return nil
	}
	for i, block := range m.src {
		m.newlineBeforeBlockIfNeeded(i)
		m.setListState(block)

		wasBlock := m.writeBlock(block)
		if m.err != nil {
			return m.err
		}

		// There may be simpler way to do this, but this
		// effectively trims all but one newline at the end.
		if wasBlock || i < len(m.src)-1 || len(m.src) == 1 {
			m.writeNewline()
		}
	}
	return m.err
}

func (m *markdownWriter) writeBlock(b portabletext.Block) bool {
	switch b.Type {
	case "image":
		m.writeImage(b.Image)
		return true
	case "code":
		m.writeCode(b.Code)
		return true
	}

	if !b.HasText() {
		return false
	}

	clear(m.markDefs)

	size := len(b.Children)

	// Resizable slice of marks.
	if cap(m.marksOpen) < size {
		m.marksOpen = make([][]string, 0, size+10)
		m.marksClose = make([][]string, 0, size+10)
	}

	m.marksOpen = m.marksOpen[:size]
	m.marksClose = m.marksClose[:size]
	m.clearMarks()

	for i, c := range b.Children {
		atStart := i == 0
		atEnd := i == size-1
		var next, prev portabletext.Child
		if !atEnd {
			next = b.Children[i+1]
		}
		if !atStart {
			prev = b.Children[i-1]
		}

		for _, mark := range c.Marks {
			inPrev, inNext := in(mark, prev.Marks), in(mark, next.Marks)
			if !inPrev {
				m.marksOpen[i] = append(m.marksOpen[i], mark)
			}
			if !inNext {
				m.marksClose[i] = append(m.marksClose[i], mark)
			}

		}

	}

	for _, md := range b.MarkDefs {
		m.markDefs[md.Key] = md
	}

	m.writeStyle(b)
	m.writeList(b)
	for i := range b.Children {
		m.writeChild(i, b.Level, b.Children)
		if m.err != nil {
			return false
		}
	}

	return true
}

func (m *markdownWriter) writeChild(i, level int, children []portabletext.Child) {
	c := children[i]
	s := c.Text
	// In the Portable Text editor, a common case is that e.g. bold text also covers whitespace
	// on either side. This works fine when rendering HTML, but not so well in Markdown.
	// E.g. this is not valid Markdown: _** italic and bold. **_
	i1 := strings.IndexFunc(s, func(r rune) bool {
		return r != ' '
	})

	i2 := strings.LastIndexFunc(s, func(r rune) bool {
		return r != ' '
	})

	var s1, s2, s3 string
	if i1 > 0 {
		s1 = s[:i1]
	}
	if i2 > 0 {
		s3 = s[i2+1:]
	}
	if i1 > 0 && i2 > 0 {
		s2 = s[i1 : i2+1]
	} else if i1 > 0 {
		s2 = s[i1:]
	} else if i2 > 0 {
		s2 = s[:i2+1]
	} else {
		s2 = s
	}

	if s1 != "" {
		m.writeString(s1)
	}
	m.writeMarks(true, i)
	m.writeString(m.indentText(s2, level))
	m.writeMarks(false, i)
	if s3 != "" {
		m.writeString(s3)
	}
}

func (m *markdownWriter) writeCode(c portabletext.Code) {
	m.writeString("```")
	m.writeString(c.Language)
	if c.Filename != "" {
		// Add as markdown attribute.
		m.writeString(fmt.Sprintf(" {filename=%q}", c.Filename))
	}
	m.writeNewline()
	m.writeString(c.Code)
	m.writeNewline()
	m.writeString("```")
}

func (m *markdownWriter) writeImage(img portabletext.Image) {
	m.writeNewline()
	m.writeString(fmt.Sprintf("![%s](%s)", img.Asset.AltText, img.Asset.URL))
}

func (m *markdownWriter) writeIndent(level int) {
	m.writeString(m.indent(level))
}

func (m *markdownWriter) writeList(b portabletext.Block) {
	if b.ListItem == "" {
		return
	}
	m.writeIndent(b.Level)

	switch b.ListItem {
	case "bullet":
		m.writeString("* ")
	case "number":
		m.writeString("1. ") // Auto numbering is provided by the editor/renderer.
	case "square":
		m.writeString("- ")
	default:
		m.writeString("* ")
	}
}

func (m *markdownWriter) writeMark(start bool, s string) {
	switch s {
	case "strong":
		m.writeString("**")
	case "code":
		m.writeString("`")
	case "em":
		m.writeString("_")
	case "underline":
		// Not supported in Markdown.
		m.writeString("")
	case "strike-through":
		m.writeString("~~")
	default:
		// Try to find the markDef.
		md, found := m.markDefs[s]
		if !found {
			return
		}
		switch md.Type {
		case "link":
			if start {
				m.writeString("[")
			} else {
				m.writeString("](")
				m.writeString(md.Href)
				m.writeString(")")
			}
		default:
			panic("unsupported mark type: " + md.Type)
		}

	}
}

func (m *markdownWriter) writeMarks(start bool, i int) {
	if start {
		for _, mark := range m.marksOpen[i] {
			m.writeMark(start, mark)
		}
	} else {
		// Reverse order.
		for j := len(m.marksClose[i]) - 1; j >= 0; j-- {
			m.writeMark(start, m.marksClose[i][j])
		}
	}
}

func (m *markdownWriter) writeNewline() {
	m.dst.Write(newlineb)
}

func (m *markdownWriter) writeString(s string) {
	_, err := io.WriteString(m.dst, s)
	if err != nil {
		m.err = err
	}
}

func (m *markdownWriter) writeStyle(b portabletext.Block) {
	style := b.Style

	switch style {
	case "normal":
	case "":
	// Do nothing.
	case "blockquote":
		m.writeString("> ")
	case "code":
		m.writeString("```")
	case "pre":
		m.writeString("```")
	default:
		if b.ListItem != "" {
			return
		}
		switch style {
		case "h1":
			m.writeString("# ")
		case "h2":
			m.writeString("## ")
		case "h3":
			m.writeString("### ")
		case "h4":
			m.writeString("#### ")
		case "h5":
			m.writeString("##### ")
		case "h6":
			m.writeString("###### ")
		default:
			return
		}
	}
}

func (m *markdownWriter) clear() {
	m.dst = nil
	m.src = nil
	m.err = nil
	m.currentListItem = ""

	clear(m.markDefs)
	m.clearMarks()
}

func (m *markdownWriter) clearMarks() {
	for i := range m.marksOpen {
		m.marksOpen[i] = m.marksOpen[i][:0]
	}
	for i := range m.marksClose {
		m.marksClose[i] = m.marksClose[i][:0]
	}
}

func in(e string, s []string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

var writerPool = sync.Pool{
	New: func() any {
		return &markdownWriter{
			markDefs: make(map[string]portabletext.MarkDef),
		}
	},
}

func getWriter() *markdownWriter {
	return writerPool.Get().(*markdownWriter)
}

func putWriter(w *markdownWriter) {
	w.clear()
	writerPool.Put(w)
}
