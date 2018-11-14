package taskpaper

import (
	"bytes"
	"strings"
	"testing"
)

func TestUnmarshalMarshal(t *testing.T) {
	tests := []string{
		`- A`,
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
	for _, test := range tests {
		in := []byte(test)
		doc, err := Unmarshal(in)
		if err != nil {
			t.Fatal(err)
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
