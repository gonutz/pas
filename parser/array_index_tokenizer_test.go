package parser

import "testing"

func TestArrayIndex(t *testing.T) {
	checkTokens(t,
		`1..20`,
		tok(tokenWord, "1"),
		tok('.', "."),
		tok('.', "."),
		tok(tokenWord, "20"),
		tok(tokenEOF, ""),
	)
}
