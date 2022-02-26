package parser

import (
	"testing"

	"github.com/akm/pas/ast"
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
							&ast.Record{
								Name: "T1",
								Members: []ast.ClassMember{
									&ast.Variable{Name: "A", Type: "Integer"},
								},
							},
							&ast.Record{
								Name: "T2",
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		})
}
