package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var Para = "I am a person, you can call me on <PhoneNumber> or [Email] or you can {Lexy} [[Email] [InvalidTag]"

var Tags = map[string]string{
	"PhoneNumber": "2024",
	"Email":       "steve@jobs.com",
	"Lexy":        "True",
}

type Type int8

const (
	None Type = iota
	TypeSquareBracket
	TypeSquigglyBracket
	TypeLessThanMoreThan
)

func TagReplacerViaLexer(v string, t map[string]string) string {
	var InTagType Type = None
	var JustEnteredTag = false
	var CurrentTagName strings.Builder
	var OutString strings.Builder

	for _, c := range v {
		if InTagType == None && (c == '{' || c == '<' || c == '[') {
			JustEnteredTag = true
			CurrentTagName.Reset()
			switch c {
			case '{':
				InTagType = TypeSquigglyBracket
			case '[':
				InTagType = TypeSquareBracket
			default:
				InTagType = TypeLessThanMoreThan
			}
		} else if JustEnteredTag {
			JustEnteredTag = false
			// we're probably escaping here
			if c == '{' && InTagType == TypeSquigglyBracket {
				OutString.WriteRune('{')
				InTagType = None
			} else if c == '[' && InTagType == TypeSquareBracket {
				OutString.WriteRune('[')
				InTagType = None
			} else if c == '<' && InTagType == TypeLessThanMoreThan {
				OutString.WriteRune('<')
				InTagType = None
			} else {
				CurrentTagName.WriteRune(c)
			}
		} else if (InTagType == TypeSquigglyBracket && c == '}') || (InTagType == TypeSquareBracket && c == ']') || (InTagType == TypeLessThanMoreThan && c == '>') {
			InTagType = None
			// check what the tag is
			val, ok := t[CurrentTagName.String()]
			if ok {
				OutString.WriteString(val)
			} else {
				if c == '}' {
					OutString.WriteString("{ERROR:")
					OutString.WriteString(CurrentTagName.String())
					OutString.WriteString("}")
				} else if c == ']' {
					OutString.WriteString("[ERROR:")
					OutString.WriteString(CurrentTagName.String())
					OutString.WriteString("]")
				} else {
					OutString.WriteString("<ERROR:")
					OutString.WriteString(CurrentTagName.String())
					OutString.WriteString(">")
				}
			}
		} else {
			if InTagType == None {
				OutString.WriteRune(c)
			} else {
				CurrentTagName.WriteRune(c)
			}
		}
	}

	return OutString.String()
}

func TagReplacerViaRegex(v string, t map[string]string) string {
	patterns := []struct {
		regex   *regexp.Regexp
		wrapper string
	}{
		{regexp.MustCompile(`\{\{|\{([^}]+)\}`), "{%s}"},
		{regexp.MustCompile(`\[\[|\[([^\]]+)\]`), "[%s]"},
		{regexp.MustCompile(`<<|<([^>]+)>`), "<%s>"},
	}

	result := v
	for _, p := range patterns {
		result = p.regex.ReplaceAllStringFunc(result, func(match string) string {
			// Check if this is an escaped tag
			if len(match) == 2 {
				return match[:1] // Return single bracket for escaped tags
			}

			// Extract tag name by removing the brackets
			tagName := match[1 : len(match)-1]

			if val, ok := t[tagName]; ok {
				return val
			}
			return fmt.Sprintf(p.wrapper, "ERROR:"+tagName)
		})
	}

	return result
}

func BenchmarkRegexMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = TagReplacerViaRegex(Para, Tags)
	}
}

func TestRegexMethod(b *testing.T) {
	fmt.Println(TagReplacerViaRegex(Para, Tags))
}

func TestLexerMethod(t *testing.T) {
	fmt.Println(TagReplacerViaLexer(Para, Tags))
}

func BenchmarkLexerMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = TagReplacerViaLexer(Para, Tags)
	}
}
