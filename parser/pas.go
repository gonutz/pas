package parser

import (
	"github.com/akm/delparser/ast"
)

func ParseString(code string) (*ast.File, error) {
	return new([]rune(code)).parseFile()
}
