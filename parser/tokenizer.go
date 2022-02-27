package parser

import (
	"unicode"
)

func newTokenizer(code []rune) tokenizer {
	return tokenizer{
		code: code,
		line: 1,
		col:  1,
	}
}

type tokenizer struct {
	code []rune
	cur  int
	line int
	col  int
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
func isDigitFirst(r rune) bool {
	return isDigit(r) || r == '-' // || r == '+' // Don't allow + for now
}

func (t *tokenizer) next() token {
	haveType := tokenIllegal
	start := t.cur
	line, col := t.line, t.col

	r := t.currentRune()
	switch r {
	case 0:
		return token{
			tokenType: tokenEOF,
			line:      line,
			col:       col,
		}
	case ';', ':', '.', ',', '=', '(', ')', '[', ']':
		t.nextRune()
		haveType = tokenType(r)
	case '{':
		for {
			r := t.nextRune()
			if r == '}' || r == 0 {
				break
			}
		}
		t.nextRune()
		haveType = tokenComment
	case '/':
		if t.nextRune() == '/' {
			for {
				r := t.nextRune()
				if r == '\n' || r == 0 {
					break
				}
			}
			t.nextRune()
			haveType = tokenComment
		}
	default:
		if unicode.IsSpace(r) {
			for unicode.IsSpace(t.nextRune()) {
			}
			haveType = tokenWhiteSpace
		} else if r == '_' || unicode.IsLetter(r) {
			word := func(r rune) bool {
				return r == '_' || unicode.IsLetter(r) || isDigit(r)
			}
			for word(t.nextRune()) {
			}
			haveType = tokenWord
		} else if isDigitFirst(r) {
			hasDot := false
			for {
				r := t.nextRune()
				if isDigit(r) {
					// OK
				} else if r == '.' {
					if t.seek(1) == '.' {
						break
					} else {
						if hasDot {
							break
						}
						hasDot = true
					}
				} else {
					break
				}
			}
			if hasDot {
				haveType = tokenReal
			} else {
				haveType = tokenInt
			}
		} else {
			t.nextRune()
		}
	}

	return token{
		tokenType: haveType,
		text:      string(t.code[start:t.cur]),
		line:      line,
		col:       col,
	}
}

func (t *tokenizer) currentRune() rune {
	if t.cur < len(t.code) {
		return t.code[t.cur]
	}
	return 0
}

func (t *tokenizer) nextRune() rune {
	if t.cur < len(t.code) {
		if t.code[t.cur] == '\n' {
			t.line++
			t.col = 1
		} else {
			t.col++
		}
		t.cur++
	}
	return t.currentRune()
}

func (t *tokenizer) seek(n int) rune {
	if t.cur+n < len(t.code) {
		return t.code[t.cur+n]
	}
	return 0
}
