package parser

import (
	"testing"

	"github.com/akm/delparser/ast"
)

// https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E6%A7%8B%E9%80%A0%E5%8C%96%E5%9E%8B%EF%BC%88Delphi%EF%BC%89#.E9.9D.99.E7.9A.84.E9.85.8D.E5.88.97
func Test2StaticArray(t *testing.T) {
	matrixIndexTypes := []ast.IndexType{
		&ast.NumRange{Low: 1, High: 10},
		&ast.NumRange{Low: 1, High: 50},
	}

	packedMatrixIndexTypes := []ast.IndexType{
		&ast.NamedIndexType{Packed: true, Name: "Boolean"},
		&ast.NumRange{Packed: true, Low: 1, High: 10},
		&ast.NamedIndexType{Packed: true, Name: "TShoeSize"},
	}

	parseFile(t, `
  unit U;
  interface
  type
    TNumbers = packed array [1..100] of Real;
    TMyArray = array [1..100] of Char;
    TMatrix1 = array[1..10] of array[1..50] of Real;
	TMatrix2 = array[1..10, 1..50] of Real;
	TMyPackedArray1 = packed array[Boolean, 1..10, TShoeSize] of Integer;
	TMyPackedArray2 = packed array[Boolean] of packed array[1..10] of packed array[TShoeSize] of Integer;
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
							ast.NewArray("TNumbers", "Real", []ast.IndexType{&ast.NumRange{Packed: true, Low: 1, High: 100}}),
							ast.NewArray("TMyArray", "Char", []ast.IndexType{&ast.NumRange{Low: 1, High: 100}}),
							ast.NewArray("TMatrix1", "Real", matrixIndexTypes),
							ast.NewArray("TMatrix2", "Real", matrixIndexTypes),
							ast.NewArray("TMyPackedArray1", "Integer", packedMatrixIndexTypes),
							ast.NewArray("TMyPackedArray2", "Integer", packedMatrixIndexTypes),
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		})
}
