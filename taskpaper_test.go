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
	tests := []string{
		`- A`,
		`My Note`,
		`My Project:`,
		strings.TrimSpace(`
- A
	- B
`),
		strings.TrimSpace(`
- A
		- B
	- C
`),
		strings.TrimSpace(`
- A
	- B
	- C
		- D
			- E
		- F
- G
- H
`),
	}

	for i, test := range tests {
		in := []byte(test)
		doc, err := Unmarshal(in)
		if err != nil {
			t.Error(err)
			continue
		}
		jsonDoc, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			t.Error(err)
			continue
		}

		fixture := fmt.Sprintf("%d.golden", i+1)
		if err := gc.GoldenFixture(jsonDoc, fixture); err != nil {
			t.Error(err)
			continue
		}

		out, err := Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(in, out) {
			t.Errorf("\nGot:\n%s\n\nWant:\n%s\n", out, in)
		}
	}
}
