package parser

import "testing"

func TestFloatValue(t *testing.T) {
	checkTokens(t,
		`3.5`,
		tok(tokenWord, "3.5"),
		tok(tokenEOF, ""),
	)
	checkTokens(t,
		`1.2.3.5`,
		tok(tokenWord, "1.2"),
		tok('.', "."),
		tok(tokenWord, "3.5"),
		tok(tokenEOF, ""),
	)
}
