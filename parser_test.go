package pas_test

import (
	"strings"
	"testing"

	"github.com/gonutz/check"
	"github.com/gonutz/pas"
)

func TestParseEmptyUnit(t *testing.T) {
	parseFile(t,
		`unit U;
	interface
	implementation
	end.`,
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{Kind: pas.InterfaceSection},
				{Kind: pas.ImplementationSection},
			},
		})
}

func TestUnitWithDotsInName(t *testing.T) {
	parseFile(t,
		`unit U.V.W;
	interface
	implementation
	end.`,
		&pas.File{
			Kind: pas.Unit,
			Name: "U.V.W",
			Sections: []pas.FileSection{
				{Kind: pas.InterfaceSection},
				{Kind: pas.ImplementationSection},
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
  end.
`,
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Uses: []string{
						"CustomUnit",
						"System.Math",
						"Vcl.Graphics.Splines",
					},
				},
				{
					Kind: pas.ImplementationSection,
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
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Blocks: []pas.FileSectionBlock{
						pas.TypeBlock{
							pas.Class{
								Name: "C",
							},
						},
					},
				},
				{Kind: pas.ImplementationSection},
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
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Blocks: []pas.FileSectionBlock{
						pas.TypeBlock{
							pas.Class{
								Name:         "G",
								SuperClasses: []string{"A", "B.C", "D.E.F"},
							},
						},
					},
				},
				{Kind: pas.ImplementationSection},
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
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Blocks: []pas.FileSectionBlock{
						pas.TypeBlock{
							pas.Class{
								Name: "C",
								Sections: []pas.ClassSection{
									{Members: []pas.ClassMember{
										pas.Variable{Name: "A", Type: "Integer"},
										pas.Variable{Name: "B", Type: "C.D"},
									}},
								},
							},
						},
					},
				},
				{Kind: pas.ImplementationSection},
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
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Blocks: []pas.FileSectionBlock{
						pas.TypeBlock{
							pas.Class{
								Name: "C", Sections: []pas.ClassSection{
									{Members: []pas.ClassMember{
										pas.Function{Name: "A"},
										pas.Function{Name: "B"},
										pas.Function{Name: "C",
											Parameters: []pas.Parameter{
												{
													Names: []string{"D"},
													Type:  "Integer",
												},
											},
										},
										pas.Function{Name: "E",
											Parameters: []pas.Parameter{
												{
													Names: []string{"F", "G"},
													Type:  "Integer",
												},
											},
										},
										pas.Function{Name: "H",
											Parameters: []pas.Parameter{
												{
													Names: []string{"I"},
													Type:  "Integer",
												},
												{
													Names: []string{"J"},
													Type:  "string",
												},
											},
										},
										pas.Function{Name: "A", Returns: "Integer"},
										pas.Function{Name: "B", Returns: "string"},
										pas.Function{Name: "C", Returns: "Pointer",
											Parameters: []pas.Parameter{
												{
													Names: []string{"D"},
													Type:  "Integer",
												},
											},
										},
										pas.Function{Name: "E", Returns: "Cardinal",
											Parameters: []pas.Parameter{
												{
													Names: []string{"F", "G"},
													Type:  "Integer",
												},
											},
										},
										pas.Function{Name: "H", Returns: "Vcl.TForm",
											Parameters: []pas.Parameter{
												{
													Names: []string{"I"},
													Type:  "Integer",
												},
												{
													Names: []string{"J"},
													Type:  "K.L",
												},
											},
										},
										pas.Function{Name: "P",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: pas.Var,
												},
											},
										},
										pas.Function{Name: "P",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: pas.Const,
												},
											},
										},
										pas.Function{Name: "P",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: pas.ConstRef,
												},
											},
										},
										pas.Function{Name: "P",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: pas.RefConst,
												},
											},
										},
										pas.Function{Name: "P",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"I"},
													Type:      "Integer",
													Qualifier: pas.Out,
												},
											},
										},
										pas.Function{Name: "NoType",
											Parameters: []pas.Parameter{
												{
													Names:     []string{"P"},
													Type:      "",
													Qualifier: pas.Const,
												},
											},
										},
									}},
								},
							},
						},
					},
				},
				{Kind: pas.ImplementationSection},
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
		&pas.File{
			Kind: pas.Unit,
			Name: "U",
			Sections: []pas.FileSection{
				{
					Kind: pas.InterfaceSection,
					Blocks: []pas.FileSectionBlock{
						pas.TypeBlock{
							pas.Class{
								Name: "C", Sections: []pas.ClassSection{
									{
										Visibility: pas.DefaultPublished,
										Members: []pas.ClassMember{
											pas.Variable{Name: "A", Type: "Integer"},
										},
									},
									{
										Visibility: pas.Public,
										Members: []pas.ClassMember{
											pas.Variable{Name: "B", Type: "Integer"},
										},
									},
									{
										Visibility: pas.Private,
										Members: []pas.ClassMember{
											pas.Variable{Name: "C", Type: "Integer"},
										},
									},
									{
										Visibility: pas.Protected,
										Members: []pas.ClassMember{
											pas.Variable{Name: "D", Type: "Integer"},
										},
									},
									{
										Visibility: pas.Published,
										Members: []pas.ClassMember{
											pas.Variable{Name: "E", Type: "Integer"},
										},
									},
								},
							},
						},
					},
				},
				{Kind: pas.ImplementationSection},
			},
		},
	)
}

func parseFile(t *testing.T, code string, want *pas.File) {
	t.Helper()
	code = strings.Replace(code, "\n", "\r\n", -1)
	f, err := pas.ParseString(code)
	if err != nil {
		t.Fatal(err)
	}
	check.Eq(t, f, want)
}
