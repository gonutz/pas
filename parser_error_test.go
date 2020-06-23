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
		`keyword "unit" expected but was end of file at 1:1`,
	)
	parseError(t,
		"unit",
		`unit name expected but was end of file at 1:5`,
	)
	parseError(t,
		"unit U",
		`token ";" expected but was end of file at 1:7`,
	)
	parseError(t,
		"unit U;",
		`keyword "interface" expected but was end of file at 1:8`,
	)
	parseError(t,
		"unit U;interface",
		`keyword "implementation" expected but was end of file at 1:17`,
	)
	parseError(t,
		"unit U;interface implementation",
		`keyword "end" expected but was end of file at 1:32`,
	)
	parseError(t,
		"unit U;interface implementation end",
		`token "." expected but was end of file at 1:36`,
	)
}

func TestIncompleteUses(t *testing.T) {
	// Valid code that we break at different points:
	//
	//     unit U;interface uses GR32, System.StrUtils; implementation end.
	parseError(t,
		"unit U;interface uses GR32, System.StrUtils implementation end.",
		`token ";" expected but was word "implementation" at 1:45`,
	)
	parseError(t,
		"unit U;interface uses GR32, System.; implementation end.",
		`uses clause expected but was token ";" at 1:36`,
	)
	parseError(t,
		"unit U;interface uses GR32, ; implementation end.",
		`uses clause expected but was token ";" at 1:29`,
	)
	parseError(t,
		"unit U;interface uses GR32 System.StrUtils; implementation end.",
		`token ";" expected but was word "System" at 1:28`,
	)
	parseError(t,
		"unit U;interface uses , System.StrUtils; implementation end.",
		`uses clause expected but was token "," at 1:23`,
	)
	parseError(t,
		"unit U;interface uses ; implementation end.",
		`uses clause expected but was token ";" at 1:23`,
	)
	parseError(t,
		"unit U;interface uses implementation end.",
		`token ";" expected but was word "end" at 1:38`,
	)
}

func TestIncompleteTypeBlock(t *testing.T) {
	// Valid code that we break at different points:
	//
	//     unit U;interface type C=class(A,B) end; implementation end.
	parseError(t,
		"unit U;interface type C=class(A,B) end implementation end.",
		`token ";" expected but was word "implementation" at 1:40`,
	)
	parseError(t,
		"unit U;interface type C=class(A,B) ; implementation end.",
		`keyword "end" expected but was token ";" at 1:36`,
	)
	parseError(t,
		"unit U;interface type C=class(A,B end; implementation end.",
		`token ")" expected but was word "end" at 1:35`,
	)
	parseError(t,
		"unit U;interface type C=class(A,) end; implementation end.",
		`parent interface name expected but was token ")" at 1:33`,
	)
	parseError(t,
		"unit U;interface type C=class(A B) end; implementation end.",
		`token ")" expected but was word "B" at 1:33`,
	)
	parseError(t,
		"unit U;interface type C=class(,B) end; implementation end.",
		`parent class name expected but was token "," at 1:31`,
	)
	parseError(t,
		"unit U;interface type C=A,B) end; implementation end.",
		`keyword "class" expected but was word "A" at 1:25`,
	)
	parseError(t,
		"unit U;interface type C class(A,B) end; implementation end.",
		`token "=" expected but was word "class" at 1:25`,
	)
	parseError(t,
		"unit U;interface type =class(A,B) end; implementation end.",
		`type name expected but was token "=" at 1:23`,
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
