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
										&ast.Method{Function: ast.Function{Name: "B"}},
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
										&ast.Method{Function: ast.Function{Name: "B", Returns: "string"}},
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
