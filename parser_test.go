package pas_test

import (
	"strings"
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/pas"
)

func TestParseEmptyUnit(t *testing.T) {
	parseFile(t,
		`unit U;
	interface
	implementation
	end.`,
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
		})
}

func parseFile(t *testing.T, code string, want *pas.File) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	f, err := pas.ParseString(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, f, want)
}
