package parser

import "testing"

func TestTokenizeString(t *testing.T) {
	checkTokens(t,
		`''`,
		tok(tokenString, "''"),
		tok(tokenEOF, ""),
	)

	checkTokens(t,
		`'foo'`,
		tok(tokenString, "'foo'"),
		tok(tokenEOF, ""),
	)

	checkTokens(t,
		`'escaped \'single quotes\''`,
		tok(tokenString, "'escaped \\'single quotes\\''"),
		tok(tokenEOF, ""),
	)
}
