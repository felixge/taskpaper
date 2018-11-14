# taskpaper

This library implements the [TaskPaper](https://www.taskpaper.com/) file format
in [Go](https://golang.org/).

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
  "Content": "",
  "Indent": 0,
  "Children": [
    {
      "Kind": "project",
      "Content": "My Project",
      "Indent": 0,
      "Children": [
        {
          "Kind": "task",
          "Content": "Task A",
          "Indent": 0,
          "Children": [
            {
              "Kind": "task",
              "Content": "Task B",
              "Indent": 0,
              "Children": null
            }
          ]
        },
        {
          "Kind": "task",
          "Content": "Task C",
          "Indent": 0,
          "Children": [
            {
              "Kind": "note",
              "Content": "Some note for Task C",
              "Indent": 0,
              "Children": null
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
