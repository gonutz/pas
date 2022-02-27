package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/delparser/ast"
)

func TestParseEmptyClass(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type C = class end;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							&ast.Class{
								Name: "C",
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

func TestParseInheritingClass(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type G = class(A, B.C, D.E.F) end;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							&ast.Class{
								Name:         "G",
								SuperClasses: []string{"A", "B.C", "D.E.F"},
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

func TestParseClassFields(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type C = class
    A: Integer;
    B: C.D;
  end;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							&ast.Class{
								Name: "C",
								Sections: []ast.ClassSection{
									{Members: []ast.ClassMember{
										&ast.Field{Variable: ast.Variable{Name: "A", Type: "Integer"}},
										&ast.Field{Variable: ast.Variable{Name: "B", Type: "C.D"}},
									}},
								},
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

func TestParseClassFunctions(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type C = class
    procedure A;
    procedure B();
    procedure C(D: Integer);
    procedure E(F, G: Integer);
    procedure H(I: Integer; J: string);
    function A: Integer;
    function B(): string;
    function C(D: Integer): Pointer;
    function E(F, G: Integer): Cardinal;
    function H(I: Integer; J: K.L): Vcl.TForm;
    procedure P(var I: Integer);
    procedure P(const I: Integer);
    procedure P(const [Ref] I: Integer);
    procedure P([Ref] const I: Integer);
    procedure P(out I: Integer);
    procedure NoType(const P);
  end;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							&ast.Class{
								Name: "C", Sections: []ast.ClassSection{
									{Members: []ast.ClassMember{
										&ast.Method{Function: ast.Function{Name: "A"}},
										&ast.Method{Function: ast.Function{Name: "B", Parameters: ast.Parameters{}}},
										&ast.Method{Function: ast.Function{Name: "C",
											Parameters: ast.Parameters{
												{
													Names: []string{"D"},
													Type:  "Integer",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "E",
											Parameters: ast.Parameters{
												{
													Names: []string{"F", "G"},
													Type:  "Integer",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "H",
											Parameters: ast.Parameters{
												{
													Names: []string{"I"},
													Type:  "Integer",
												},
												{
													Names: []string{"J"},
													Type:  "string",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "A", Returns: "Integer"}},
										&ast.Method{Function: ast.Function{Name: "B", Returns: "string", Parameters: ast.Parameters{}}},
										&ast.Method{Function: ast.Function{Name: "C", Returns: "Pointer",
											Parameters: ast.Parameters{
												{
													Names: []string{"D"},
													Type:  "Integer",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "E", Returns: "Cardinal",
											Parameters: ast.Parameters{
												{
													Names: []string{"F", "G"},
													Type:  "Integer",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "H", Returns: "Vcl.TForm",
											Parameters: ast.Parameters{
												{
													Names: []string{"I"},
													Type:  "Integer",
												},
												{
													Names: []string{"J"},
													Type:  "K.L",
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "P",
											Parameters: ast.Parameters{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: ast.Var,
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "P",
											Parameters: ast.Parameters{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: ast.Const,
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "P",
											Parameters: ast.Parameters{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: ast.ConstRef,
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "P",
											Parameters: ast.Parameters{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: ast.RefConst,
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "P",
											Parameters: ast.Parameters{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: ast.Out,
												},
											},
										}},
										&ast.Method{Function: ast.Function{Name: "NoType",
											Parameters: ast.Parameters{
												{
													Names:     []string{"P"},
													Type:      "",
													Qualifier: ast.Const,
												},
											},
										}},
									}},
								},
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

func TestClassVisibilities(t *testing.T) {
	parseFile(t, `
  unit U;
  interface
  type C = class
    A: Integer;
  public
    B: Integer;
  private
    C: Integer;
  protected
    D: Integer;
  published
    E: Integer;
  end;
  implementation
  end.`,
		&ast.File{
			Kind: ast.Unit,
			Name: "U",
			Sections: []*ast.FileSection{
				{
					Kind: ast.InterfaceSection,
					Blocks: []ast.FileSectionBlock{
						ast.TypeBlock{
							&ast.Class{
								Name: "C", Sections: []ast.ClassSection{
									{
										Visibility: ast.DefaultPublished,
										Members: []ast.ClassMember{
											&ast.Field{Variable: ast.Variable{Name: "A", Type: "Integer"}},
										},
									},
									{
										Visibility: ast.Public,
										Members: []ast.ClassMember{
											&ast.Field{Variable: ast.Variable{Name: "B", Type: "Integer"}},
										},
									},
									{
										Visibility: ast.Private,
										Members: []ast.ClassMember{
											&ast.Field{Variable: ast.Variable{Name: "C", Type: "Integer"}},
										},
									},
									{
										Visibility: ast.Protected,
										Members: []ast.ClassMember{
											&ast.Field{Variable: ast.Variable{Name: "D", Type: "Integer"}},
										},
									},
									{
										Visibility: ast.Published,
										Members: []ast.ClassMember{
											&ast.Field{Variable: ast.Variable{Name: "E", Type: "Integer"}},
										},
									},
								},
							},
						},
					},
				},
				{Kind: ast.ImplementationSection},
			},
		},
	)
}

// https://docwiki.embarcadero.com/RADStudio/Alexandria/en/Properties_(Delphi)
func TestParseClassProperties(t *testing.T) {
	type pattern struct {
		code string
		prop *ast.Property
	}

	patterns := []*pattern{
		{
			"property Color: TColor read GetColor write SetColor;",
			&ast.Property{Variable: ast.Variable{Name: "Color", Type: "TColor"}, Reader: "GetColor", Writer: "SetColor"},
		},
		{
			"property Objects[Index: Integer]: TObject read GetObject write SetObject;",
			&ast.Property{Variable: ast.Variable{Name: "Objects", Type: "TObject"}, Reader: "GetObject", Writer: "SetObject", Indexes: []*ast.Parameter{
				{Names: []string{"Index"}, Type: "Integer"},
			}},
		},
		{
			"property Pixels[X, Y: Integer]: TColor read GetPixel write SetPixel;",
			&ast.Property{Variable: ast.Variable{Name: "Pixels", Type: "TColor"}, Reader: "GetPixel", Writer: "SetPixel", Indexes: []*ast.Parameter{
				{Names: []string{"X", "Y"}, Type: "Integer"},
			}},
		},
		{
			"property Values[const Name: string]: string read GetValue write SetValue;",
			&ast.Property{Variable: ast.Variable{Name: "Values", Type: "string"}, Reader: "GetValue", Writer: "SetValue", Indexes: []*ast.Parameter{
				{Names: []string{"Name"}, Type: "string", Qualifier: ast.Const},
			}},
		},
		{
			"property Left:   Longint index 0 read GetCoordinate write SetCoordinate;",
			&ast.Property{Variable: ast.Variable{Name: "Left", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 0},
		},
		{
			"property Top:    Longint index 1 read GetCoordinate write SetCoordinate;",
			&ast.Property{Variable: ast.Variable{Name: "Top", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 1},
		},
		{
			"property Right:  Longint index 2 read GetCoordinate write SetCoordinate;",
			&ast.Property{Variable: ast.Variable{Name: "Right", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 2},
		},
		{
			"property Bottom: Longint index 3 read GetCoordinate write SetCoordinate;",
			&ast.Property{Variable: ast.Variable{Name: "Bottom", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 3},
		},
		{
			"property Coordinates[Index: Integer]: Longint read GetCoordinate write SetCoordinate;",
			&ast.Property{Variable: ast.Variable{Name: "Coordinates", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Indexes: []*ast.Parameter{
				{Names: []string{"Index"}, Type: "Integer"},
			}},
		},
		{
			"property Name: TComponentName read FName write SetName stored False;",
			&ast.Property{Variable: ast.Variable{Name: "Name", Type: "TComponentName"}, Reader: "FName", Writer: "SetName", Stored: "False"},
		},
		{
			"property Tag: Longint read FTag write FTag default 0;",
			&ast.Property{Variable: ast.Variable{Name: "Tag", Type: "Longint"}, Reader: "FTag", Writer: "FTag", Default: "0"},
		},
		{
			"property Name: string read FName write FName default 'User1';",
			&ast.Property{Variable: ast.Variable{Name: "Name", Type: "string"}, Reader: "FName", Writer: "FName", Default: "'User1'"},
		},
		{
			"class property Red: Integer read FRed write FRed;",
			&ast.Property{Variable: ast.Variable{Name: "Red", Type: "Integer"}, Reader: "FRed", Writer: "FRed", Class: true},
		},
	}

	unitTemplate := `unit U;
	interface
	type C = class
	%s
	end;
	implementation
	end.`

	for _, ptn := range patterns {
		t.Run(ptn.code, func(t *testing.T) {
			parseFile(t, fmt.Sprintf(unitTemplate, ptn.code),
				&ast.File{
					Kind: ast.Unit,
					Name: "U",
					Sections: []*ast.FileSection{
						{
							Kind: ast.InterfaceSection,
							Blocks: []ast.FileSectionBlock{
								ast.TypeBlock{
									&ast.Class{
										Name: "C", Sections: []ast.ClassSection{
											{Members: []ast.ClassMember{
												ptn.prop,
											}},
										},
									},
								},
							},
						},
						{Kind: ast.ImplementationSection},
					},
				},
			)
		})
	}

	t.Run("all in one", func(t *testing.T) {
		codes := make([]string, len(patterns))
		props := make([]ast.ClassMember, len(patterns))
		for i, ptn := range patterns {
			codes[i] = ptn.code
			props[i] = ptn.prop
		}
		parseFile(t, fmt.Sprintf(unitTemplate, strings.Join(codes, "\n")),
			&ast.File{
				Kind: ast.Unit,
				Name: "U",
				Sections: []*ast.FileSection{
					{
						Kind: ast.InterfaceSection,
						Blocks: []ast.FileSectionBlock{
							ast.TypeBlock{
								&ast.Class{
									Name: "C", Sections: []ast.ClassSection{
										{Members: props},
									},
								},
							},
						},
					},
					{Kind: ast.ImplementationSection},
				},
			},
		)
	})
}
