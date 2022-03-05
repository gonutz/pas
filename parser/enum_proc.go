package parser

import (
	"github.com/akm/delparser/ast"
)

func enumProc(enum *ast.EnumExpr) func(p *parser) error {
	return func(p *parser) error {
		return p.startEndToken('(', ')', func() error {
			for {

				member, err := p.identifier("enum member name")
				if err != nil {
					return err
				}
				m := ast.EnumMember{Name: member}
				if p.sees('=') {
					if err := p.eat('='); err != nil {
						return err
					}
					t, err := p.take(tokenInt)
					if err != nil {
						return err
					}
					m.Value = t.text
				}
				enum.Members = append(enum.Members, m)
				if p.sees(',') {
					if err := p.eat(','); err != nil {
						return err
					}
				} else {
					break
				}
			}
			return nil
		})
	}
}
