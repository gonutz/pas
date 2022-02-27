package parser

import (
	"strconv"

	"github.com/akm/pas/ast"
	"github.com/pkg/errors"
)

func propertyProc(res *ast.Property) func(*parser) error {
	return func(p *parser) error {
		if p.sees('[') {
			err := p.startEndToken('[', ']', func() error {
				parameters, err := p.parseParameters(']')
				if err != nil {
					return err
				}
				res.Indexes = parameters
				return nil
			})
			if err != nil {
				return err
			}
		}
		if err := p.eat(':'); err != nil {
			return err
		}
		typ, err := p.identifier("property type name")
		if err != nil {
			return err
		}
		res.Type = typ
		for !p.sees(';') {
			if p.seesWord("index") {
				if err := p.eatWord("index"); err != nil {
					return err
				}
				token := p.nextToken()
				if token.tokenType != tokenInt {
					return errors.Errorf("expected int, got %+v", token)
				}
				index, err := strconv.Atoi(token.text)
				if err != nil {
					return err
				}
				res.Index = index
			} else if p.seesWord("read") {
				if err := p.eatWord("read"); err != nil {
					return err
				}
				reader, err := p.identifier("property reader name")
				if err != nil {
					return err
				}
				res.Reader = reader
			} else if p.seesWord("write") {
				if err := p.eatWord("write"); err != nil {
					return err
				}
				writer, err := p.identifier("property writer name")
				if err != nil {
					return err
				}
				res.Writer = writer
			} else if p.seesWord("default") {
				if err := p.eatWord("default"); err != nil {
					return err
				}
				t := p.nextToken()
				switch t.tokenType {
				case tokenWord, tokenInt, tokenReal:
				// OK
				default:
					return errors.Errorf("expected property default value, got %+v", t)
				}
				res.Default = t.text
			} else if p.seesWord("stored") {
				if err := p.eatWord("stored"); err != nil {
					return err
				}
				stored, err := p.identifier("property stored value")
				if err != nil {
					return err
				}
				res.Stored = stored
			} else {
				return errors.Errorf("expected property modifier, got %+v", p.peekToken())
			}
		}
		if err := p.eat(';'); err != nil {
			return err
		}
		return nil
	}
}
