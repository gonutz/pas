package parser

import (
	"github.com/akm/pas/ast"
)

func ParseString(code string) (*ast.File, error) {
	return new([]rune(code)).parseFile()
}
