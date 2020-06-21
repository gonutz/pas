package pas_test

import (
	"strings"
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/pas"
)

func TestUnitMustHaveInterfaceAndImplementationSections(t *testing.T) {
	parseError(t,
		"",
		`keyword unit expected but was end of file`,
	)
	parseError(t,
		"unit",
		`unit name expected but was end of file`,
	)
	parseError(t,
		"unit U",
		`token ";" expected but was end of file`,
	)
	parseError(t,
		"unit U;",
		`keyword interface expected but was end of file`,
	)
	parseError(t,
		"unit U;interface",
		`keyword implementation expected but was end of file`,
	)
	parseError(t,
		"unit U;interface implementation",
		`keyword end expected but was end of file`,
	)
	parseError(t,
		"unit U;interface implementation end",
		`token "." expected but was end of file`,
	)
}

func parseError(t *testing.T, code, wantMessage string) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	_, err := pas.ParseString(code)
	if err == nil {
		t.Fatal("error expected")
	}
	check.Eq(t, err.Error(), wantMessage)
}
