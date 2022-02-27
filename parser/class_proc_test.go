package parser

import (
	"testing"

	"github.com/akm/pas/ast"
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
	parseFile(t, `
  unit U;
  interface
  type C = class
  property Color: TColor read GetColor write SetColor;
  property Objects[Index: Integer]: TObject read GetObject write SetObject;
  property Pixels[X, Y: Integer]: TColor read GetPixel write SetPixel;
  property Values[const Name: string]: string read GetValue write SetValue;
  property Left:   Longint index 0 read GetCoordinate write SetCoordinate;
  property Top:    Longint index 1 read GetCoordinate write SetCoordinate;
  property Right:  Longint index 2 read GetCoordinate write SetCoordinate;
  property Bottom: Longint index 3 read GetCoordinate write SetCoordinate;
  property Coordinates[Index: Integer]: Longint read GetCoordinate write SetCoordinate;
  property Name: TComponentName read FName write SetName stored False;
  property Tag: Longint read FTag write FTag default 0;
  class property Red: Integer read FRed write FRed;
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
										&ast.Property{Variable: ast.Variable{Name: "Color", Type: "TColor"}, Reader: "GetColor", Writer: "SetColor"},
										&ast.Property{Variable: ast.Variable{Name: "Objects", Type: "TObject"}, Reader: "GetObject", Writer: "SetObject", Indexes: []*ast.Parameter{
											{Names: []string{"Index"}, Type: "Integer"},
										}},
										&ast.Property{Variable: ast.Variable{Name: "Pixels", Type: "TColor"}, Reader: "GetPixel", Writer: "SetPixel", Indexes: []*ast.Parameter{
											{Names: []string{"X", "Y"}, Type: "Integer"},
										}},
										&ast.Property{Variable: ast.Variable{Name: "Values", Type: "string"}, Reader: "GetValue", Writer: "SetValue", Indexes: []*ast.Parameter{
											{Names: []string{"Name"}, Type: "string", Qualifier: ast.Const},
										}},
										&ast.Property{Variable: ast.Variable{Name: "Left", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 0},
										&ast.Property{Variable: ast.Variable{Name: "Top", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 1},
										&ast.Property{Variable: ast.Variable{Name: "Right", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 2},
										&ast.Property{Variable: ast.Variable{Name: "Bottom", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Index: 3},
										&ast.Property{Variable: ast.Variable{Name: "Coordinates", Type: "Longint"}, Reader: "GetCoordinate", Writer: "SetCoordinate", Indexes: []*ast.Parameter{
											{Names: []string{"Index"}, Type: "Integer"},
										}},
										&ast.Property{Variable: ast.Variable{Name: "Name", Type: "TComponentName"}, Reader: "FName", Writer: "SetName", Stored: "False"},
										&ast.Property{Variable: ast.Variable{Name: "Tag", Type: "Longint"}, Reader: "FTag", Writer: "FTag", Default: "0"},
										&ast.Property{Variable: ast.Variable{Name: "Red", Type: "Integer"}, Reader: "FRed", Writer: "FRed", Class: true},
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
