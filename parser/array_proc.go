package parser

import (
	"regexp"
	"strconv"

	"github.com/akm/delparser/ast"
	"github.com/pkg/errors"
)

var ptnDigits = regexp.MustCompile(`[0-9]+`)

func arrayProc(array *ast.Array) func(p *parser) error {
	return func(p *parser) error {
		for {
			packed := false
			if p.seesWord("packed") {
				if err := p.eatWord("packed"); err != nil {
					return err
				}
				packed = true
			}
			if err := p.eatWord("array"); err != nil {
				return err
			}
			if p.sees('[') {
				if err := p.eat('['); err != nil {
					return err
				}
				for {
					var indexType ast.IndexType
					token := p.nextToken()
					if ptnDigits.MatchString(token.text) {
						low, err := strconv.Atoi(token.text)
						if err != nil {
							return errors.Wrapf(err, "invalid array index token %+v", token)
						}
						if err := p.eats('.', '.'); err != nil {
							return err
						}
						token := p.nextToken()
						high, err := strconv.Atoi(token.text)
						if err != nil {
							return errors.Wrapf(err, "invalid array index token %+v", token)
						}
						indexType = &ast.NumRange{Packed: packed, Low: low, High: high}
					} else {
						indexType = &ast.NamedIndexType{Packed: packed, Name: token.text}
					}
					array.IndexTypes = append(array.IndexTypes, indexType)

					if p.sees(',') {
						p.eat(',')
					} else if p.sees(']') {
						break
					}
				}

				if err := p.eat(']'); err != nil {
					return err
				}
			}
			if err := p.eatWord("of"); err != nil {
				return err
			}
			if !p.seesWords("packed", "array") {
				break
			}
		}
		typ, err := p.identifier("array type name")
		if err != nil {
			return err
		}
		array.Type = typ
		if err := p.eat(';'); err != nil {
			return err
		}
		return nil
	}
}
