package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/pas/ast"
)

// https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E3%83%91%E3%83%A9%E3%83%A1%E3%83%BC%E3%82%BF%EF%BC%88Delphi%EF%BC%89
func TestParameters(t *testing.T) {
	type pattern struct {
		code     string
		expected *ast.Function
	}

	patterns := []*pattern{
		{
			"function Power(X: Real; Y: Integer): Real;",
			&ast.Function{Name: "Power", Returns: "Real", Parameters: ast.Parameters{
				{Names: []string{"X"}, Type: "Real"},
				{Names: []string{"Y"}, Type: "Integer"},
			}},
		},
		{
			"procedure A(X, Y: Real);",
			&ast.Function{Name: "A", Parameters: ast.Parameters{
				{Names: []string{"X", "Y"}, Type: "Real"},
			}},
		},
		{
			"procedure B(var S: string; X: Integer);",
			&ast.Function{Name: "B", Parameters: ast.Parameters{
				{Names: []string{"S"}, Type: "string", Qualifier: ast.Var},
				{Names: []string{"X"}, Type: "Integer"},
			}},
		},
		{
			"procedure C(HWnd: Integer; Text, Caption: PChar; Flags: Integer);",
			&ast.Function{Name: "C", Parameters: ast.Parameters{
				{Names: []string{"HWnd"}, Type: "Integer"},
				{Names: []string{"Text", "Caption"}, Type: "PChar"},
				{Names: []string{"Flags"}, Type: "Integer"},
			}},
		},
		{
			"procedure D(const P; I: Integer);",
			&ast.Function{Name: "D", Parameters: ast.Parameters{
				{Names: []string{"P"}, Type: "", Qualifier: ast.Const},
				{Names: []string{"I"}, Type: "Integer"},
			}},
		},
		{
			"procedure UpdateRecords;",
			&ast.Function{Name: "UpdateRecords"},
		},
		{
			"function DoubleByValue(X: Integer): Integer;",
			&ast.Function{Name: "DoubleByValue", Returns: "Integer", Parameters: ast.Parameters{
				{Names: []string{"X"}, Type: "Integer"},
			}},
		},
		{
			"function DoubleByRef(var X: Integer): Integer;",
			&ast.Function{Name: "DoubleByRef", Returns: "Integer", Parameters: ast.Parameters{
				{Names: []string{"X"}, Type: "Integer", Qualifier: ast.Var},
			}},
		},
		{
			"function CompareStr(const S1, S2: string): Integer;",
			&ast.Function{Name: "CompareStr", Returns: "Integer", Parameters: ast.Parameters{
				{Names: []string{"S1", "S2"}, Type: "string", Qualifier: ast.Const},
			}},
		},
		{
			"function FunctionName(const [Ref] parameter1: Class1Name; [Ref] const parameter2: Class2Name);",
			&ast.Function{Name: "FunctionName", Parameters: ast.Parameters{
				{Names: []string{"parameter1"}, Type: "Class1Name", Qualifier: ast.ConstRef},
				{Names: []string{"parameter2"}, Type: "Class2Name", Qualifier: ast.RefConst},
			}},
		},
		{
			"procedure GetInfo(out Info: SomeRecordType);",
			&ast.Function{Name: "GetInfo", Parameters: ast.Parameters{
				{Names: []string{"Info"}, Type: "SomeRecordType", Qualifier: ast.Out},
			}},
		},
		{
			"procedure TakeAnything(const C);",
			&ast.Function{Name: "TakeAnything", Parameters: ast.Parameters{
				{Names: []string{"C"}, Type: "", Qualifier: ast.Const},
			}},
		},
		{
			"function Equal(var Source, Dest; Size: Integer): Boolean;",
			&ast.Function{Name: "Equal", Returns: "Boolean", Parameters: ast.Parameters{
				{Names: []string{"Source", "Dest"}, Type: "", Qualifier: ast.Var},
				{Names: []string{"Size"}, Type: "Integer"},
			}},
		},
		{
			"procedure Check(S: OpenString);",
			&ast.Function{Name: "Check", Parameters: ast.Parameters{{Names: []string{"S"}, Type: "OpenString"}}},
		},
		{
			"procedure Sort(A: TDigits);",
			&ast.Function{Name: "Sort", Parameters: ast.Parameters{{Names: []string{"A"}, Type: "TDigits"}}},
		},
		{
			"function Find(A: array of Char): Integer; // Open Array Parameters",
			&ast.Function{Name: "Find", Returns: "Integer", Parameters: ast.Parameters{
				{Names: []string{"A"}, Type: "Char", OpenArray: true},
			}},
		},
		{
			"function MakeStr(const Args: array of const): string; // Variant Open Array Parameters",
			&ast.Function{Name: "MakeStr", Returns: "string", Parameters: ast.Parameters{
				{Names: []string{"Args"}, Type: "const", OpenArray: true, Qualifier: ast.Const},
			}},
		},
		{
			"procedure FillArray(A: array of Integer; Value: Integer = 0);",
			&ast.Function{Name: "FillArray", Parameters: ast.Parameters{
				{Names: []string{"A"}, Type: "Integer", OpenArray: true},
				{Names: []string{"Value"}, Type: "Integer", DefaultValue: "0"},
			}},
		},
		{
			"function MyFunction1(X: Real = 3.5; Y: Real = 3.5): Real;",
			&ast.Function{Name: "MyFunction1", Returns: "Real", Parameters: ast.Parameters{
				{Names: []string{"X"}, Type: "Real", DefaultValue: "3.5"},
				{Names: []string{"Y"}, Type: "Real", DefaultValue: "3.5"},
			}},
		},
		{
			"function MyFunction2(X: Real = -3.5; Y: Real = -3.5): Real;",
			&ast.Function{Name: "MyFunction2", Returns: "Real", Parameters: ast.Parameters{
				{Names: []string{"X"}, Type: "Real", DefaultValue: "-3.5"},
				{Names: []string{"Y"}, Type: "Real", DefaultValue: "-3.5"},
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
			code := fmt.Sprintf(unitTemplate, ptn.code)
			parseFile(t, code, &ast.File{
				Kind: ast.Unit,
				Name: "U",
				Sections: []*ast.FileSection{
					{
						Kind: ast.InterfaceSection,
						Blocks: []ast.FileSectionBlock{
							ptn.expected,
						},
					},
					{Kind: ast.ImplementationSection},
				},
			})
		})
	}

	t.Run("all in one", func(t *testing.T) {
		codes := make([]string, len(patterns))
		functions := make([]ast.FileSectionBlock, len(patterns))
		for i, code := range patterns {
			codes[i] = code.code
			functions[i] = code.expected
		}
		parseFile(t, fmt.Sprintf(unitTemplate, strings.Join(codes, "\n")),
			&ast.File{
				Kind: ast.Unit,
				Name: "U",
				Sections: []*ast.FileSection{
					{
						Kind:   ast.InterfaceSection,
						Blocks: functions,
					},
					{Kind: ast.ImplementationSection},
				},
			},
		)
	})
}
