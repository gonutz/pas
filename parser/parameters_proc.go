package parser

import "github.com/akm/pas/ast"

func parametersProc(dest *ast.Parameters, endToken tokenType) func(*parser) error {
	return func(p *parser) error {
		res := ast.Parameters{}
		for p.sees(tokenWord) || p.sees('[') {
			param := &ast.Parameter{}
			if p.seesWord("var") {
				if err := p.eatWord("var"); err != nil {
					return err
				}
				param.Qualifier = ast.Var
			} else if p.seesWord("const") {
				if err := p.eatWord("const"); err != nil {
					return err
				}
				param.Qualifier = ast.Const
				if p.sees('[') {
					err := p.startEndToken('[', ']', func() (err error) {
						return p.eatWord("ref")
					})
					if err != nil {
						return err
					}
					param.Qualifier = ast.ConstRef
				}
			} else if p.seesWord("out") {
				if err := p.eatWord("out"); err != nil {
					return err
				}
				param.Qualifier = ast.Out
			} else if p.sees('[') {
				err := p.startEndToken('[', ']', func() (err error) {
					return p.eatWord("ref")
				})
				if err != nil {
					return err
				}
				if err := p.eatWord("const"); err != nil {
					return err
				}
				param.Qualifier = ast.RefConst
			}

			names, err := p.parseSeparatedString(',', "parameter name")
			if err != nil {
				return err
			}
			param.Names = names

			if p.sees(':') {
				if err := p.eat(':'); err != nil {
					return err
				}
				pt, err := p.qualifiedIdentifier("parameter type")
				if err != nil {
					return err
				}
				param.Type = pt
			}
			res = append(res, param)
			if p.sees(';') {
				p.eat(';')
			} else if p.sees(endToken) {
				break
			} else {
				if err := p.eat(','); err != nil {
					return err
				}
				break // The last parameter is not followed by a ';'.
			}
		}
		*dest = res
		return nil
	}
}
