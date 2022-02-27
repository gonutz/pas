package parser

import (
	"testing"

	"github.com/akm/pas/ast"
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

func TestParseVarBlock(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  VAR
    I: Integer;
    S: string;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.VarBlock{
							{Name: "I", Type: "Integer"},
							{Name: "S", Type: "string"},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

func TestParseTwoVarBlocks(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  var I: Integer;
  var S: string;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.VarBlock{{Name: "I", Type: "Integer"}},
						ast.VarBlock{{Name: "S", Type: "string"}},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
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
