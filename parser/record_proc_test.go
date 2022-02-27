package parser

import (
	"testing"

	"github.com/akm/delparser/ast"
)

func Test2RecordsInATypeBlock(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type
    T1 = record
	  A: Integer;
    end;
    T2 = record
    end;
  implementation
  {$R *.dfm}
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							ast.NewRecord("T1",
								&ast.Field{Variable: *ast.NewVariable("A", "Integer")},
							),
							ast.NewRecord("T2"),
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		})
}
