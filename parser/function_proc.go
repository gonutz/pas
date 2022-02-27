package parser

import "github.com/akm/delparser/ast"

func functionProc(f *ast.Function) func(*parser) error {
	return func(p *parser) error {
		if p.sees('(') {
			err := p.startEndToken('(', ')', func() (err error) {
				parameters, err := p.parseParameters(')')
				if err != nil {
					return err
				}
				f.Parameters = parameters
				return nil
			})
			if err != nil {
				return err
			}
		}
		if p.sees(':') {
			if err := p.eat(':'); err != nil {
				return err
			}
			rt, err := p.qualifiedIdentifier("return type")
			if err != nil {
				return err
			}
			f.Returns = rt
		}
		if err := p.eat(';'); err != nil {
			return err
		}
		return nil
	}
}
