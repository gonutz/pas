package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		`field name expected but was token ";" at 1:36`,
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
		`keyword "record" expected but was word "A" at 1:25`,
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

func TestIncompleteVarInClass(t *testing.T) {
	// Valid code that we break at different points:
	//
	//     unit U;interface type C=class A:Integer; end; implementation end.
	parseError(t,
		"unit U;interface type C=class A:Integer end; implementation end.",
		`token ";" expected but was word "end" at 1:41`,
	)
	parseError(t,
		"unit U;interface type C=class A:; end; implementation end.",
		`type name expected but was token ";" at 1:33`,
	)
	parseError(t,
		"unit U;interface type C=class A Integer; end; implementation end.",
		`token ":" expected but was word "Integer" at 1:33`,
	)
	parseError(t,
		"unit U;interface type C=class :Integer; end; implementation end.",
		`field name expected but was token ":" at 1:31`,
	)
}

func TestIncompleteClassFunctions(t *testing.T) {
	parseError(t,
		"unit U;interface type C=class procedure( end; implementation end.",
		`function name expected but was token "(" at 1:40`,
	)
	// The following test parses a procedure named "end" and expectes the class
	// field "implementation" to be followed by a type, e.g. ": Integer".
	// TODO This is not a good message for this error.
	parseError(t,
		"unit U;interface type C=class procedure end; implementation end.",
		`token ":" expected but was word "end" at 1:61`,
	)
	// The same as above happens here.
	parseError(t,
		"unit U;interface type C=class function A: end; implementation end.",
		`token ":" expected but was word "end" at 1:63`,
	)
}

func parseError(t *testing.T, code, wantMessage string) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	_, err := ParseString(code)
	if err == nil {
		t.Fatal("error expected")
	}
	assert.Equal(t, wantMessage, err.Error())
}
