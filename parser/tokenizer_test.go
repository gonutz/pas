package parser

import "testing"

func TestTokenize(t *testing.T) {
	checkTokens(t,
		`unit, U=()[];: end.%`,
		tok(tokenWord, "unit"),
		tok(',', ","),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "U"),
		tok('=', "="),
		tok('(', "("),
		tok(')', ")"),
		tok('[', "["),
		tok(']', "]"),
		tok(';', ";"),
		tok(':', ":"),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "end"),
		tok('.', "."),
		tok(tokenIllegal, "%"),
		tok(tokenEOF, ""),
	)
}

func TestTokenizeComments(t *testing.T) {
	checkTokens(t,
		`{this is a
comment} {another}//and a line comment
//plus a line comment just before EOF`,
		tok(tokenComment, "{this is a\ncomment}"),
		tok(tokenWhiteSpace, " "),
		tok(tokenComment, "{another}"),
		tok(tokenComment, "//and a line comment\n"),
		tok(tokenComment, "//plus a line comment just before EOF"),
		tok(tokenEOF, ""),
	)
}
