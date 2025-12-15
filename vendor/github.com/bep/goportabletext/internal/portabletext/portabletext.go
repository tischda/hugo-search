// Copyright 2025 Bjørn Erik Pedersen
// SPDX-License-Identifier: MIT

package portabletext

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"unicode"

	"github.com/mitchellh/mapstructure"
)

type Type string

const (
	TypeBlock Type = "block"
	TypeSpan  Type = "span"
)

type Blocks []Block

type BaseBlock struct {
	// Key is a unique identifier for the block.
	Key string `json:"_key" mapstructure:"_key"`
	// The type makes it possible for a serializer to parse the contents of the block.
	Type Type `json:"_type" mapstructure:"_type"`
}

// A Block is a top-level structure in a Portable Text array.
type Block struct {
	BaseBlock `mapstructure:",squash"`

	// Union type.
	// A block can be either a text block or an image block.
	// The type field is used to determine which type of block it is.
	Text  `mapstructure:",squash"`
	Image `mapstructure:",squash"`
	Code  `mapstructure:",squash"`
}

type Style string

type Text struct {
	// Children is an array of spans or custom inline types that is contained within a block.
	Children []Child `json:"children" mapstructure:"children"`

	// MarkDefs definitions is an array of objects with a key, type and some data.
	// Mark definitions are tied to spans by adding the referring _key in the marks array.
	MarkDefs []MarkDef `json:"markDefs" mapstructure:"markDefs"`

	// Style typically describes a visual property for the whole block.
	Style Style `json:"style,omitempty" mapstructure:"style"`

	// Level is used to express visual nesting and hierarchical structures between blocks in the array.
	Level int `json:"level,omitempty" mapstructure:"level"`

	// A block can be given the property listItem with a value that describes which kind of list it is.
	// Typically bullet, number, square and so on.
	// The list position is derived from the position the block has in the array and surrounding list items on the same level.
	ListItem string `json:"listItem,omitempty" mapstructure:"listItem"`
}

func (b Block) HasText() bool {
	for _, child := range b.Children {
		if child.Text != "" {
			return true
		}
	}
	return false
}

type Image struct {
	Asset Asset `json:"asset" mapstructure:"asset"`
}

type Asset struct {
	ID          string `json:"_id" mapstructure:"_id"`
	AltText     string `json:"altText" mapstructure:"altText"`
	Description string `json:"description" mapstructure:"description"`
	Metadata    struct {
		Dimensions struct {
			AspectRatio float64 `json:"aspectRatio" mapstructure:"aspectRatio"`
			Height      int     `json:"height" mapstructure:"height"`
			Width       int     `json:"width" mapstructure:"width"`
		} `json:"dimensions"`
	} `json:"metadata"`
	Path  string `json:"path" mapstructure:"path"`
	Title string `json:"title" mapstructure:"title"`
	URL   string `json:"url" mapstructure:"url"`
}

type Code struct {
	Code     string `json:"code" mapstructure:"code"`
	Filename string `json:"filename" mapstructure:"filename"`
	Language string `json:"language" mapstructure:"language"`
}

// A Child is a span or custom inline type that is contained within a block.
// A span is the standard way to express inline text within a block
type Child struct {
	// The type makes it possible for a serializer to parse the contents of the block.
	Type Type `json:"_type" mapstructure:"_type"`

	// Marks are how we mark up inline text with additional data/features.
	//  Marks comes in two forms: Decorators and Annotations.
	// Decorators are marks as simple string values, while Annotations are keys to a data structure.
	// marks is therefore either an array of string values, or keys, corresponding to markDefs, with the same _key.
	Marks []string `json:"marks" mapstructure:"marks"`

	// The text contents of the span.
	Text string `json:"text" mapstructure:"text"`
}

type MarkDef struct {
	Key  string `json:"_key" mapstructure:"_key"`
	Type Type   `json:"_type" mapstructure:"_type"`
	Href string `json:"href" mapstructure:"href"`
}

// The Parse function reads a Portable Text block or array from the src.
// The src aan be either a io.Reader or a ... TODO1
func Parse(src any) (Blocks, error) {
	switch v := src.(type) {
	case io.Reader:
		return parseReader(v)
	case map[string]any:
		return parseMap(v)
	case []any:
		return parseSlice(v)
	default:
		panic(fmt.Sprintf("unsupported type %T", src))
	}
}

func parseMap(m map[string]any) (Blocks, error) {
	var block Block
	if err := mapstructure.Decode(m, &block); err != nil {
		return nil, err
	}
	return Blocks{block}, nil
}

func parseSlice(s []any) (Blocks, error) {
	var blocks Blocks
	if err := mapstructure.Decode(s, &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

func parseReader(r io.Reader) (Blocks, error) {
	// Check if r represents an array of blocks or a single block.
	// We need to buffer the input to determine this.
	buff := bufio.NewReader(r)

	var firstNonSpace rune

	for {
		c, _, err := buff.ReadRune()
		if err != nil {
			return nil, err
		}
		if !unicode.IsSpace(c) {
			firstNonSpace = c
			buff.UnreadRune()
			break
		}
	}

	switch firstNonSpace {
	case '[':
		// Parse the array of blocks.
		var blocks Blocks
		if err := json.NewDecoder(buff).Decode(&blocks); err != nil {
			return nil, err
		}
		return blocks, nil
	case '{':
		// Parse a single block.
		var block Block
		if err := json.NewDecoder(buff).Decode(&block); err != nil {
			return nil, err
		}
		return Blocks{block}, nil
	default:
		return nil, nil
	}
}
