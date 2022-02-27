package parser

import (
	"testing"

	"github.com/akm/pas/ast"
)

// https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E3%83%91%E3%83%A9%E3%83%A1%E3%83%BC%E3%82%BF%EF%BC%88Delphi%EF%BC%89
func TestParameters(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  function Power(X: Real; Y: Integer): Real;
  procedure A(X, Y: Real);
  procedure B(var S: string; X: Integer);
  procedure C(HWnd: Integer; Text, Caption: PChar; Flags: Integer);
  procedure D(const P; I: Integer);
  procedure UpdateRecords;
  function DoubleByValue(X: Integer): Integer;
  function DoubleByRef(var X: Integer): Integer;
  function CompareStr(const S1, S2: string): Integer;
  function FunctionName(const [Ref] parameter1: Class1Name; [Ref] const parameter2: Class2Name);
  procedure GetInfo(out Info: SomeRecordType);
  procedure TakeAnything(const C);
  function Equal(var Source, Dest; Size: Integer): Boolean;
  // procedure Check(S: string[20]); // syntax error
  procedure Check(S: OpenString);
  procedure Sort(A: TDigits);
  function Find(A: array of Char): Integer; // Open Array Parameters
  function MakeStr(const Args: array of const): string; // Variant Open Array Parameters
//    procedure FillArray(A: array of Integer; Value: Integer = 0);
//   function MyFunction(X: Real = 3.5; Y: Real = 3.5): Real;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						&ast.Function{Name: "Power", Returns: "Real", Parameters: ast.Parameters{
							{Names: []string{"X"}, Type: "Real"},
							{Names: []string{"Y"}, Type: "Integer"},
						}},
						&ast.Function{Name: "A", Parameters: ast.Parameters{
							{Names: []string{"X", "Y"}, Type: "Real"},
						}},
						&ast.Function{Name: "B", Parameters: ast.Parameters{
							{Names: []string{"S"}, Type: "string", Qualifier: ast.Var},
							{Names: []string{"X"}, Type: "Integer"},
						}},
						&ast.Function{Name: "C", Parameters: ast.Parameters{
							{Names: []string{"HWnd"}, Type: "Integer"},
							{Names: []string{"Text", "Caption"}, Type: "PChar"},
							{Names: []string{"Flags"}, Type: "Integer"},
						}},
						&ast.Function{Name: "D", Parameters: ast.Parameters{
							{Names: []string{"P"}, Type: "", Qualifier: ast.Const},
							{Names: []string{"I"}, Type: "Integer"},
						}},
						&ast.Function{Name: "UpdateRecords"},
						&ast.Function{Name: "DoubleByValue", Returns: "Integer", Parameters: ast.Parameters{
							{Names: []string{"X"}, Type: "Integer"},
						}},
						&ast.Function{Name: "DoubleByRef", Returns: "Integer", Parameters: ast.Parameters{
							{Names: []string{"X"}, Type: "Integer", Qualifier: ast.Var},
						}},
						&ast.Function{Name: "CompareStr", Returns: "Integer", Parameters: ast.Parameters{
							{Names: []string{"S1", "S2"}, Type: "string", Qualifier: ast.Const},
						}},
						&ast.Function{Name: "FunctionName", Parameters: ast.Parameters{
							{Names: []string{"parameter1"}, Type: "Class1Name", Qualifier: ast.ConstRef},
							{Names: []string{"parameter2"}, Type: "Class2Name", Qualifier: ast.RefConst},
						}},
						&ast.Function{Name: "GetInfo", Parameters: ast.Parameters{
							{Names: []string{"Info"}, Type: "SomeRecordType", Qualifier: ast.Out},
						}},
						&ast.Function{Name: "TakeAnything", Parameters: ast.Parameters{
							{Names: []string{"C"}, Type: "", Qualifier: ast.Const},
						}},
						&ast.Function{Name: "Equal", Returns: "Boolean", Parameters: ast.Parameters{
							{Names: []string{"Source", "Dest"}, Type: "", Qualifier: ast.Var},
							{Names: []string{"Size"}, Type: "Integer"},
						}},
						&ast.Function{Name: "Check", Parameters: ast.Parameters{{Names: []string{"S"}, Type: "OpenString"}}},
						&ast.Function{Name: "Sort", Parameters: ast.Parameters{{Names: []string{"A"}, Type: "TDigits"}}},
						&ast.Function{Name: "Find", Returns: "Integer", Parameters: ast.Parameters{
							{Names: []string{"A"}, Type: "Char", OpenArray: true},
						}},
						&ast.Function{Name: "MakeStr", Returns: "string", Parameters: ast.Parameters{
							{Names: []string{"Args"}, Type: "const", OpenArray: true, Qualifier: ast.Const},
						}},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}
