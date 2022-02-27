package parser

import "testing"

func TestFloatValue(t *testing.T) {
	checkTokens(t,
		`3.5`,
		tok(tokenReal, "3.5"),
		tok(tokenEOF, ""),
	)
	checkTokens(t,
		`1.2.3`,
		tok(tokenReal, "1.2"),
		tok('.', "."),
		tok(tokenInt, "3"),
		tok(tokenEOF, ""),
	)
	checkTokens(t,
		`1.2.3.5`,
		tok(tokenReal, "1.2"),
		tok('.', "."),
		tok(tokenReal, "3.5"),
		tok(tokenEOF, ""),
	)
	checkTokens(t,
		`-3.5`,
		tok(tokenReal, "-3.5"),
		tok(tokenEOF, ""),
	)
	checkTokens(t,
		`-5`,
		tok(tokenInt, "-5"),
		tok(tokenEOF, ""),
	)
}
