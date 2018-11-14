# taskpaper

This library implements the [TaskPaper](https://www.taskpaper.com/) file format
in [Go](https://golang.org/).

[![GoDoc](https://godoc.org/github.com/felixge/taskpaper?status.svg)](https://godoc.org/github.com/felixge/taskpaper)
[![Build Status](https://travis-ci.org/felixge/taskpaper.svg?branch=master)](https://travis-ci.org/felixge/taskpaper)

# Usage

Below is a short example program that parses TaskPaper document and shows
the resulting data structure:

```go
doc, _ := taskpaper.Unmarshal([]byte(strings.TrimSpace(`
My Project:
- Task A
  - Task B
- Task C
  Some note for Task C
`)))
data, _ := json.MarshalIndent(doc, "", "  ")
fmt.Printf("%s\n", data)
```
```json

{
  "Kind": "document",
  "Children": [
    {
      "Kind": "project",
      "Content": "My Project",
      "Children": [
        {
          "Kind": "task",
          "Content": "Task A",
          "Children": [
            {
              "Kind": "task",
              "Content": "Task B"
            }
          ]
        },
        {
          "Kind": "task",
          "Content": "Task C",
          "Children": [
            {
              "Kind": "note",
              "Content": "Some note for Task C"
            }
          ]
        }
      ]
    }
  ]
}

```

# File Format

Since there is only a [reference implementation](https://www.taskpaper.com/),
but no specification for the file format, I decided to define the TaskPaper
[ABNF](https://en.wikipedia.org/wiki/Augmented_Backus%E2%80%93Naur_form) that
is implemented by this library as follows:

```
document = *(item CR) *1item ; the last item has no CR
item     = (project / task / note)
project  = *OCTET ":"
task     = "- " *OCTET
note     = *OCTET
```
