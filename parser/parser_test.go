package parser

import (
	"fmt"
	"testing"

	"github.com/akm/delparser/ast"
)

func TestParseEmptyUnit(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{Kind: ast.InterfaceSection},
				{Kind: ast.ImplementationSection},
			},
		})
}

func TestUnitWithDotsInName(t *testing.T) {
	parseFile(t, `
  unit U.V.W;
  interface
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U.V.W",
			Sections: []*ast.FileSection{
				{Kind: ast.InterfaceSection},
				{Kind: ast.ImplementationSection},
			},
		})
}

func TestParseUses(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  uses CustomUnit, System.Math, Vcl.Graphics.Splines;
  implementation
  uses Windows . WinAPI;
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Uses: []string{
						"CustomUnit",
						"System.Math",
						"Vcl.Graphics.Splines",
					},
				},
				{
					Kind: ast.ImplementationSection,
					Uses: []string{
						"Windows.WinAPI",
					},
				},
			},
		})
}

func TestInterfaceSection(t *testing.T) {
	type pattern struct {
		name   string
		code   string
		blocks []ast.FileSectionBlock
	}

	patterns := []pattern{
		{
			"2 declarations in a var block",
			`VAR
			I: Integer;
			S: string;`,
			[]ast.FileSectionBlock{
				ast.VarBlock{
					ast.NewVariable("I", "Integer"),
					ast.NewVariable("S", "string"),
				},
			},
		},
		{
			"2 var blocks",
			`var I: Integer;
			var S: string;`,
			[]ast.FileSectionBlock{
				ast.VarBlock{ast.NewVariable("I", "Integer")},
				ast.VarBlock{ast.NewVariable("S", "string")},
			},
		},
		{
			"var blocks and function",
			`var I: Integer;
			function Bar: string;`,
			[]ast.FileSectionBlock{
				ast.VarBlock{ast.NewVariable("I", "Integer")},
				&ast.Function{Name: "Bar", Returns: "string"},
			},
		},
		{
			"var blocks and procedure",
			`var I: Integer;
			procedure B(var S: string; X: Integer);`,
			[]ast.FileSectionBlock{
				ast.VarBlock{ast.NewVariable("I", "Integer")},
				&ast.Function{Name: "B", Parameters: ast.Parameters{
					{Names: []string{"S"}, Type: "string", Qualifier: ast.Var},
					{Names: []string{"X"}, Type: "Integer"},
				}},
			},
		},
	}

	unitTemplate := `unit U;
	interface
	%s
	implementation
	end.`

	for _, ptn := range patterns {
		t.Run(ptn.name, func(t *testing.T) {
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

func TestComments(t *testing.T) {
	// TODO For now we skip comments and throw them away. Figure out a way to
	// include them in the tree, maybe as a parallel structure with references
	// to the tree.
	parseFile(t, `
  unit U;
  interface
  implementation
  {$R *.dfm}
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{Kind: ast.InterfaceSection},
				{Kind: ast.ImplementationSection},
			},
		})
}
