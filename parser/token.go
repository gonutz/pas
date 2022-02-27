package parser

import (
	"fmt"
	"unicode/utf8"
)

type token struct {
	tokenType tokenType
	text      string
	// line and col both start at 1.
	line, col int
}

// tokenType is a rune because single characters are used directly as their
// token type, e.g. ',' '+' or ':'.
type tokenType rune

const (
	tokenIllegal    tokenType = -1
	tokenEOF        tokenType = 0
	tokenWord       tokenType = 256
	tokenWhiteSpace tokenType = 257
	tokenComment    tokenType = 258
	tokenInt        tokenType = 259
	tokenReal       tokenType = 260
	tokenString     tokenType = 261
)

func (t token) String() string {
	if t.tokenType == tokenComment {
		text := t.text
		const max = 20
		if utf8.RuneCountInString(text) > max {
			text = string([]rune(text)[:max]) + "..."
		}
		return fmt.Sprintf("%v %q at %d:%d", t.tokenType, text, t.line, t.col)
	}
	if string(t.tokenType) == t.text || t.text == "" {
		return fmt.Sprintf("%v at %d:%d", t.tokenType, t.line, t.col)
	}
	return fmt.Sprintf("%v %q at %d:%d", t.tokenType, t.text, t.line, t.col)
}

func (t tokenType) String() string {
	switch t {
	case tokenIllegal:
		return "illegal token"
	case tokenEOF:
		return "end of file"
	case tokenWord:
		return "word"
	case tokenWhiteSpace:
		return "white space"
	case tokenComment:
		return "comment"
	case tokenInt:
		return "int"
	case tokenReal:
		return "real"
	case tokenString:
		return "string"
	default:
		if 0 <= t && t <= 127 {
			return fmt.Sprintf("token %q", string(t))
		}
		return fmt.Sprintf("token %q (%d)", string(t), int(t))
	}
}
