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
				ast.NewVariable("I", "Integer"),
			}},
		},
		{
			"var X, Y: Real;",
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("X", "Y", "Real"),
			}},
		},
		{
			"var Digit: 0..9;",
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("Digit", "0..9"),
			}},
		},
		{
			"var I: Integer = 7;",
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("I", "Integer").WithDefault("7"),
			}},
		},
		{
			// https://docwiki.embarcadero.com/RADStudio/Alexandria/en/String_Types_(Delphi)#Short_Strings
			"var MyString: string[100];",
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("MyString", "string").WithLength(100),
			}},
		},
		{
			`var Str: string[32];
			StrLen: Byte absolute Str;`,
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("Str", "string").WithLength(32),
				ast.NewVariable("StrLen", "Byte").WithAbsolute("Str"),
			}},
		},
		{
			"var Checks: Array [1..3] of Boolean;",
			[]ast.FileSectionBlock{ast.VarBlock{
				ast.NewVariable("Checks", "Array [1..3] of Boolean"),
			}},
		},
		{
			"threadvar X: Integer;",
			[]ast.FileSectionBlock{ast.ThreadVarBlock{
				ast.NewVariable("X", "Integer"),
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
