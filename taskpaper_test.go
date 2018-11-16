package taskpaper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/felixge/goldy"
)

var gc = goldy.DefaultConfig()

func TestUnmarshalMarshal(t *testing.T) {
	tests := []struct {
		Name      string
		TaskPaper string
	}{
		{
			Name:      "single-task",
			TaskPaper: `- A`,
		},
		{
			Name:      "single-note",
			TaskPaper: `My Note`,
		},
		{
			Name:      "single-project",
			TaskPaper: `My Project:`,
		},
		{
			Name: "nested-task",
			TaskPaper: strings.TrimSpace(`
- A
	- B
`),
		},
		{
			Name: "nested-task-indent",
			TaskPaper: strings.TrimSpace(`
- A
		- B
	- C
`),
		},
		{
			Name: "nested-task-deep",
			TaskPaper: strings.TrimSpace(`
- A
	- B
	- C
		- D
			- E
		- F
- G
- H
`),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			in := []byte(test.TaskPaper)
			doc, err := Unmarshal(in)
			if err != nil {
				t.Error(err)
				return
			}
			jsonDoc, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				t.Error(err)
				return
			}

			fixture := fmt.Sprintf("%s.golden", test.Name)
			fixtureData := []byte(fmt.Sprintf("%s\n---\n%s",
				in,
				jsonDoc,
			))
			if err := gc.GoldenFixture(fixtureData, fixture); err != nil {
				t.Error(err)
				return
			}

			out, err := Marshal(doc)
			if err != nil {
				t.Error(err)
			} else if !bytes.Equal(in, out) {
				t.Errorf("\nGot:\n%s\n\nWant:\n%s\n", out, in)
			}
		})
	}
}
