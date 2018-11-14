package taskpaper

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Item is a line item in a TaskPaper document.
type Item struct {
	// Kind is the item kind.
	Kind Kind
	// Content is the item text, not including the "- " prefix for Kind=Task or
	// the ":" suffix for Kind=Project.
	Content string `json:",omitempty"`
	// Indent is the amount of additional tabs the original item was indented by
	// if it's more than the single tab that establishes it as a child of the
	// previous line item.
	Indent int `json:",omitempty"`
	// Children is the list of all items that are indented by at least one tab
	// below the current item.
	Children []*Item
	// Parent points to the parent item.
	Parent *Item `json:"-"`
}

// Kind specifies the type of an Item.
type Kind string

const (
	// Document is the kind of the top-level item holding the contents of a
	// TaskPaper document.
	Document Kind = "document"
	// Project is a project item, i.e. a line that ended with ":".
	Project Kind = "project"
	// Task is a task item, i.e. a line that starts with "- " after indentation.
	Task Kind = "task"
	// Note is a note item, i.e. any other line that is not a Project or Task.
	Note Kind = "note"
)

type parseState string

const (
	parseIndent     parseState = "indent"
	parseTaskPrefix parseState = "taskPrefix"
	parseContent    parseState = "content"
)

// Unmarshal parses the given data as a TaskPaper document (see ABNF in README)
// and returns an Item with Kind=Document that holds the structure of the file
// or an error.
func Unmarshal(data []byte) (*Item, error) {
	var (
		doc     = &Item{Kind: Document}
		state   parseState
		dst     *Item
		item    *Item
		content []byte
	)

	reset := func() {
		dst = doc
		content = nil
		item = &Item{Kind: Note}
		state = parseIndent
	}

	addItem := func() {
		item.Parent = dst
		item.Content = string(content)
		if item.Kind == Note && strings.HasSuffix(item.Content, ":") {
			item.Content = strings.TrimSuffix(item.Content, ":")
			item.Kind = Project
		}
		dst.Children = append(dst.Children, item)
		reset()
	}

	reset()
	for _, c := range data {
		switch state {
		case parseIndent:
			switch c {
			case '\t':
				if len(dst.Children) > 0 {
					dst = dst.Children[len(dst.Children)-1]
				} else {
					item.Indent++
				}
			case '-':
				state = parseTaskPrefix
			case '\n':
				addItem()
			default:
				content = append(content, c)
				state = parseContent
			}

		case parseTaskPrefix:
			switch c {
			case ' ':
				item.Kind = Task
				state = parseContent
			case '\n':
				addItem()
			default:
				content = append(content, '-', c)
				state = parseContent
			}

		case parseContent:
			switch c {
			case '\n':
				addItem()
			default:
				content = append(content, c)
			}
		}
	}
	addItem()

	return doc, nil
}

// Marshal converts a Item of Kind=Document back into its text form, or returns
// an error.
func Marshal(doc *Item) ([]byte, error) {
	if doc.Kind != Document {
		return nil, errors.New("Item is not a Document")
	}

	out := bytes.NewBuffer(nil)
	err := marshal(out, doc, 0)
	return out.Bytes(), err
}

func marshal(out *bytes.Buffer, item *Item, depth int) error {
	for i, child := range item.Children {
		if depth > 0 || i > 0 {
			out.WriteString("\n")
		}
		out.Write(bytes.Repeat([]byte("\t"), depth+child.Indent))
		switch child.Kind {
		case Project:
			out.WriteString(child.Content + ":")
		case Task:
			out.WriteString("- " + child.Content)
		case Note:
			out.WriteString(child.Content)
		default:
			return fmt.Errorf("Item has invalid kind: %q", child.Kind)
		}
		if err := marshal(out, child, depth+1); err != nil {
			return err
		}
	}
	return nil
}
