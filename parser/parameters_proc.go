package parser

import (
	"github.com/akm/pas/ast"
	"github.com/pkg/errors"
)

func parametersProc(dest *ast.Parameters, endToken tokenType) func(*parser) error {
	eatRef := func(p *parser) error {
		err := p.startEndToken('[', ']', func() (err error) {
			return p.eatWord("ref")
		})
		if err != nil {
			return err
		}
		return nil
	}

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
				if p.sees('[') {
					if err := eatRef(p); err != nil {
						return err
					}
					param.Qualifier = ast.ConstRef
				} else {
					param.Qualifier = ast.Const
				}
			} else if p.seesWord("out") {
				if err := p.eatWord("out"); err != nil {
					return err
				}
				param.Qualifier = ast.Out
			} else if p.sees('[') {
				if err := eatRef(p); err != nil {
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
				if pt == "array" {
					if err := p.eatWord("of"); err != nil {
						return err
					}
					pt, err := p.qualifiedIdentifier("open array parameter type")
					if err != nil {
						return err
					}
					param.OpenArray = true
					param.Type = pt
				} else {
					param.Type = pt
				}

				if p.sees('=') {
					if err := p.eat('='); err != nil {
						return err
					}
					t := p.nextToken()
					switch t.tokenType {
					case tokenWord, tokenInt, tokenReal, tokenString:
						// OK
					default:
						return errors.Errorf("expected parameter default value, got %+v", t)
					}
					param.DefaultValue = t.text
				}
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
