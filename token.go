package pas

import "fmt"

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
)

func (t token) String() string {
	return fmt.Sprintf("%v: %q at %d:%d", t.tokenType, t.text, t.line, t.col)
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
	default:
		if 0 <= t && t <= 127 {
			return fmt.Sprintf("token %q", string(t))
		}
		return fmt.Sprintf("token %q (%d)", string(t), int(t))
	}
}
