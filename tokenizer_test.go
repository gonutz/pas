package pas

import "testing"

func TestTokenize(t *testing.T) {
	checkTokens(t,
		`unit U; end.`,
		tok(tokenWord, "unit"),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "U"),
		tok(';', ";"),
		tok(tokenWhiteSpace, " "),
		tok(tokenWord, "end"),
		tok('.', "."),
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
