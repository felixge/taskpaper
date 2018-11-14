package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/felixge/taskpaper"
)

func main() {
	doc, _ := taskpaper.Unmarshal([]byte(strings.TrimSpace(`
My Project:
	- Task A
		- Task B
	- Task C
		Some note for Task C
`)))
	data, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Printf("%s\n", data)
}
