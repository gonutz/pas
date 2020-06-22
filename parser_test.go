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
			Sections: []pas.FileSection{
				{Kind: pas.InterfaceSection},
				{Kind: pas.ImplementationSection},
			},
		})
}

func TestParseUses(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  uses CustomUnit, System.Math, Vcl.Graphics.Splines;
  implementation
  uses Windows . WinAPI;
  end.
`,
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Uses: []string{
						"CustomUnit",
						"System.Math",
						"Vcl.Graphics.Splines",
					},
				},
				{
					Kind: pas.ImplementationSection,
					Uses: []string{
						"Windows.WinAPI",
					},
				},
			},
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
