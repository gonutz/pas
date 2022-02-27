package parser

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/akm/delparser/ast"
	"github.com/stretchr/testify/assert"
)

func parseFile(t *testing.T, code string, expected *ast.File) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	actual, err := ParseString(code)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	if !assert.Equal(t, expected, actual) {
		expectedJson, err := json.MarshalIndent(expected, "", "  ")
		assert.NoError(t, err)
		actualJson, err := json.MarshalIndent(actual, "", "  ")
		assert.NoError(t, err)
		assert.Equal(t, string(expectedJson), string(actualJson))
	}
}
