package parser

import (
	"strings"
	"testing"

	"github.com/akm/pas/ast"
	"github.com/stretchr/testify/assert"
)

func parseFile(t *testing.T, code string, want *ast.File) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	f, err := ParseString(code)
	if err != nil {
		t.Fatalf("%+v\n", err)
	}
	assert.Equal(t, want, f)
}
