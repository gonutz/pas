package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/akm/delparser/ast"
)

// https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E3%83%91%E3%83%A9%E3%83%A1%E3%83%BC%E3%82%BF%EF%BC%88Delphi%EF%BC%89
func TestEnum(t *testing.T) {
	type pattern struct {
		code     string
		expected ast.FileSectionBlock
	}

	cardMembers := []ast.EnumMember{
		{Name: "Club"},
		{Name: "Diamond"},
		{Name: "Heart"},
		{Name: "Spade"},
	}

	patterns := []*pattern{
		{
			"type Suit = (Club, Diamond, Heart, Spade);",
			ast.TypeBlock{
				&ast.Type{Name: "Suit", Expr: &ast.EnumExpr{Members: cardMembers}},
			},
		},
		{
			"var MyCard: (Club, Diamond, Heart, Spade);",
			ast.VarBlock{
				&ast.Variable{Names: []string{"MyCard"}, Type: &ast.EnumExpr{Members: cardMembers}},
			},
		},
		{
			"var Card1, Card2: (Club, Diamond, Heart, Spade);",
			ast.VarBlock{
				&ast.Variable{Names: []string{"Card1", "Card2"}, Type: &ast.EnumExpr{Members: cardMembers}},
			},
		},
		{
			"type SomeEnum1 = (e1, e2, e3 = 1);",
			ast.TypeBlock{
				&ast.Type{Name: "SomeEnum1", Expr: &ast.EnumExpr{Members: []ast.EnumMember{
					{Name: "e1"},
					{Name: "e2"},
					{Name: "e3", Value: "1"},
				}}},
			},
		},
		{
			"type SomeEnum2 = (e1=1, e2=2, e3 = 3);",
			ast.TypeBlock{
				&ast.Type{Name: "SomeEnum2", Expr: &ast.EnumExpr{Members: []ast.EnumMember{
					{Name: "e1", Value: "1"},
					{Name: "e2", Value: "2"},
					{Name: "e3", Value: "3"},
				}}},
			},
		},

		// {
		// 	"type Size = (Small = 5, Medium = 10, Large = Small + Medium);",
		// 	&ast.Type{Name: "Suit", Expr: ast.EnumExpr{Members: []ast.EnumMember{
		// 		{Name: "Small", Value: "5"},
		// 		{Name: "Medium", Value: "10"},
		// 		{Name: "Large", ValueExpr: {"Small", "+", "Medium"}}},
		// 	}},
		// },

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
