package pas

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

func tok(typ tokenType, text string) token {
	return token{tokenType: typ, text: text}
}

func checkTokens(t *testing.T, code string, want ...token) {
	t.Helper()
	have := tokenize(code)
	eq := len(want) == len(have)
	if eq {
		for i := range want {
			eq = eq &&
				want[i].tokenType == have[i].tokenType &&
				want[i].text == have[i].text
		}
	}
	if !eq {
		t.Error(printTokenComparison(want, have))
	}
}

func tokenize(code string) []token {
	lex := newTokenizer([]rune(code))
	var tokens []token
	for {
		t := lex.next()
		tokens = append(tokens, t)
		if t.tokenType == tokenEOF {
			break
		}
	}
	return tokens
}

func printTokenComparison(want, have []token) string {
	n := len(want)
	if len(have) > n {
		n = len(have)
	}
	s := "want <-> have"
	for i := 0; i < n; i++ {
		s += "\n"

		if i < len(want) && i < len(have) &&
			want[i].tokenType == have[i].tokenType &&
			want[i].text == have[i].text {
			s += "  "
		} else {
			s += "x "
		}

		if i < len(want) {
			s += want[i].String() + " "
		} else {
			s += "no more tokens"
		}

		s += "<->"

		if i < len(have) {
			s += " " + have[i].String()
		} else {
			s += "no more tokens"
		}
	}
	return s
}
