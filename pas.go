package pas

import "github.com/akm/pas/ast"

func ParseString(code string) (*ast.File, error) {
	return newParser([]rune(code)).parseFile()
}
