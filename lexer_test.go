package LexicalTokeniserTagReplacement_test

import (
	"fmt"
	LexicalTokeniserTagReplacement "github.com/CADawg/Lexical-Tokeniser-Tag-Replacement"
	"regexp"
	"testing"
)

var Para = "I am a person, you can call me on <PhoneNumber> or [Email] or you can {Lexy} [[Email] [InvalidTag]"

var Tags = map[string]string{
	"PhoneNumber": "2024",
	"Email":       "alan@example.com",
	"Lexy":        "True",
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

func TestRegexMethod(t *testing.T) {
	fmt.Println(TagReplacerViaRegex(Para, Tags))
	if TagReplacerViaRegex(Para, Tags) == "I am a person, you can call me on 2024 or alan@example.com or you can True [Email] [ERROR:InvalidTag]" {
	} else {
		t.Error("Failed")
	}

}

func TestLexerMethod(t *testing.T) {
	if LexicalTokeniserTagReplacement.ReplaceTagsInString(Para, Tags) == "I am a person, you can call me on 2024 or alan@example.com or you can True [Email] [ERROR:InvalidTag]" {
	} else {
		t.Error("Failed")
	}
}

func BenchmarkLexerMethod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = LexicalTokeniserTagReplacement.ReplaceTagsInString(Para, Tags)
	}
}
