# Lexical Tokeniser Tag Replacement

Replaces tags i.e. `[tag]`, `{tag}` and `<tag>` with the corresponding tag value provided in a dictionary.

## Install

```bash
go get -u github.com/CADawg/Lexical-Tokeniser-Tag-Replacement
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/CADawg/Lexical-Tokeniser-Tag-Replacement"
)

func main() {
    tags := map[string]string{
        "tag1": "value1",
        "tag2": "value2",
    }

    text := "This is a [tag1] and this is a {tag2} and this is a <tag3> [tag3] and this is an [[escaped tag]."

    replacedText := LexicalTokeniserTagReplacement.ReplaceTagsInString(text, tags)

    fmt.Println(replacedText)
    
    // prints: This is a value1 and this is a value2 and this is a <ERROR:tag3> [ERROR:tag3] and this is an [escaped tag].
    // will show <error:tag name> if the tag is not found in the tags dictionary
    // using [[, {{ or << will escape the tag and not replace it.
}
```