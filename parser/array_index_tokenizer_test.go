package parser

import "testing"

func TestArrayIndex(t *testing.T) {
	checkTokens(t,
		`1..20`,
		tok(tokenInt, "1"),
		tok('.', "."),
		tok('.', "."),
		tok(tokenInt, "20"),
		tok(tokenEOF, ""),
	)
}
