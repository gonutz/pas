package parser

import (
	"fmt"
	"testing"

	"github.com/akm/delparser/ast"
)

// https://docwiki.embarcadero.com/RADStudio/Alexandria/en/Variables_(Delphi)
func TestVar(t *testing.T) {
	type pattern struct {
		code   string
		blocks []ast.FileSectionBlock
	}

	patterns := []pattern{
		{
			"var I: Integer;",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"I"}, Type: "Integer"},
			}},
		},
		{
			"var X, Y: Real;",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"X", "Y"}, Type: "Real"},
			}},
		},
		{
			"var Digit: 0..9;",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"Digit"}, Type: "0..9"},
			}},
		},
		{
			"var I: Integer = 7;",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"I"}, Type: "Integer", Default: "7"},
			}},
		},
		{
			// https://docwiki.embarcadero.com/RADStudio/Alexandria/en/String_Types_(Delphi)#Short_Strings
			"var MyString: string[100];",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"MyString"}, Type: "string", Length: 100},
			}},
		},
		{
			`var Str: string[32];
			StrLen: Byte absolute Str;`,
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"Str"}, Type: "string", Length: 32},
				{Names: []string{"StrLen"}, Type: "Byte", Absolute: "Str"},
			}},
		},
		{
			"var Checks: Array [1..3] of Boolean;",
			[]ast.FileSectionBlock{ast.VarBlock{
				{Names: []string{"Checks"}, Type: "Boolean"},
			}},
		},
		{
			"threadvar X: Integer;",
			[]ast.FileSectionBlock{ast.ThreadVarBlock{
				{Names: []string{"X"}, Type: "Integer"},
			}},
		},
	}

	unitTemplate := `unit U;
	interface
	%s
	implementation
	end.`

	for _, ptn := range patterns {
		t.Run(ptn.code, func(t *testing.T) {
			parseFile(t, fmt.Sprintf(unitTemplate, ptn.code), &ast.File{
				Kind: ast.Unit,
				Name: "U",
				Sections: []*ast.FileSection{
					{
						Kind:   ast.InterfaceSection,
						Blocks: ptn.blocks,
					},
					{Kind: ast.ImplementationSection},
				},
			})
		})
	}
}
