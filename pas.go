package delparser

import (
	"github.com/akm/delparser/ast"
	"github.com/akm/delparser/parser"
)

func ParseString(code string) (*ast.File, error) {
	return parser.ParseString(code)
}
