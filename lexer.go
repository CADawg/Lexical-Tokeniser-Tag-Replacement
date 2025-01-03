package LexicalTokeniserTagReplacement

import "strings"

type Type int8

const (
	None Type = iota << 1
	TypeSquareBracket
	TypeSquigglyBracket
	TypeLessThanMoreThan
	//TypeAll = TypeSquigglyBracket + TypeSquareBracket + TypeLessThanMoreThan
)

func ReplaceTagsInString(v string, t map[string]string) string {
	var InTagType = None
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
