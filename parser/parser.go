package parser

import (
	"strings"

	"github.com/akm/pas/ast"
	"github.com/pkg/errors"
)

func new(code []rune) *parser {
	return &parser{tokens: newTokenizer(code)}
}

type parser struct {
	tokens tokenizer
	// isPeeking and peekingAt are a one-element queue of tokens to come. The
	// tokenizer only gives us the next token, it cannot peek so we buffer one
	// token here. See parser.nextToken and parser.peekToken.
	isPeeking bool
	peekingAt token
	file      ast.File
}

func (p *parser) parseFile() (*ast.File, error) {
	// For now only parse units until we have tests for other kinds.
	if err := p.eatWord("unit"); err != nil {
		return nil, err
	}
	p.file.Kind = ast.Unit
	unitName, err := p.qualifiedIdentifier("unit name")
	if err != nil {
		return nil, err
	}
	p.file.Name = unitName
	if err := p.eat(';'); err != nil {
		return nil, err
	}

	if err := p.eatWord("interface"); err != nil {
		return nil, err
	}
	if err := p.parseFileSection(ast.InterfaceSection); err != nil {
		return nil, err
	}

	if err := p.eatWord("implementation"); err != nil {
		return nil, err
	}
	if err := p.parseFileSection(ast.ImplementationSection); err != nil {
		return nil, err
	}

	if err := p.eatWord("end"); err != nil {
		return nil, err
	}
	if err := p.eat('.'); err != nil {
		return nil, err
	}
	return &p.file, nil
}

func (p *parser) parseFileSection(kind ast.FileSectionKind) error {
	uses, err := p.parseUses()
	if err != nil {
		return err
	}
	blocks, err := p.parseSectionBlocks()
	if err != nil {
		return err
	}
	p.file.Sections = append(p.file.Sections, &ast.FileSection{
		Kind:   kind,
		Uses:   uses,
		Blocks: blocks,
	})
	return nil
}

func (p *parser) parseSeparatedString(separator tokenType, identifierDesc string, extraIdentifierDescs ...string) ([]string, error) {
	res := []string{}
	s, err := p.qualifiedIdentifier(identifierDesc)
	if err != nil {
		return nil, err
	}
	res = append(res, s)

	if len(extraIdentifierDescs) > 0 {
		identifierDesc = extraIdentifierDescs[0]
	}
	for p.sees(separator) {
		if err := p.eat(separator); err != nil {
			return nil, err
		}
		s, err := p.qualifiedIdentifier(identifierDesc)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (p *parser) startWordEndToken(start string, tt tokenType, fn func() error) error {
	if err := p.eatWord(start); err != nil {
		return err
	}
	if err := fn(); err != nil {
		return err
	}
	if err := p.eat(tt); err != nil {
		return err
	}
	return nil
}

func (p *parser) parseUses() ([]string, error) {
	if p.seesWord("uses") {
		var uses []string
		err := p.startWordEndToken("uses", ';', func() (err error) {
			uses, err = p.parseSeparatedString(',', "uses clause")
			return
		})
		return uses, err
	} else {
		return nil, nil
	}
}

func (p *parser) parseSectionBlocks() ([]ast.FileSectionBlock, error) {
	var blocks []ast.FileSectionBlock
	for {
		if p.seesWord("type") {
			typeBlock, err := p.parseTypeBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, typeBlock)
		} else if p.seesWord("var") {
			varBlock, err := p.parseVarBlock()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, varBlock)
		} else {
			break
		}
	}
	return blocks, nil
}

func (p *parser) parseTypeBlock() (ast.TypeBlock, error) {
	if err := p.eatWord("type"); err != nil {
		return nil, err
	}
	res := ast.TypeBlock{}
	for {
		identifier, err := p.identifier("type name")
		if err != nil {
			return nil, err
		}
		if err := p.eat('='); err != nil {
			return nil, err
		}
		if p.seesWord("class") {
			class := &ast.Class{Name: identifier}
			if err := classProcessor(class)(p); err != nil {
				return nil, err
			}
			res = append(res, class)
		} else if p.seesWord("record") {
			record := &ast.Record{Name: identifier}
			if err := recordProcessor(record)(p); err != nil {
				return nil, err
			}
			res = append(res, record)
		} else if p.seesWords("packed", "array") {
			array := &ast.Array{Name: identifier}
			if err := arrayProcessor(array)(p); err != nil {
				return nil, err
			}
			res = append(res, array)
		} else {
			return nil, errors.Errorf("expected type declaration, got %+v", p.peekToken())
		}
		if p.seesWords("var", "type", "const", "implementation") {
			break
		}
	}
	return res, nil
}

func (p *parser) parseVarBlock() (ast.VarBlock, error) {
	if err := p.eatWord("var"); err != nil {
		return nil, err
	}
	var vars ast.VarBlock
	for p.sees(tokenWord) && !p.seesKeyword() {
		varDec, err := p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}
		vars = append(vars, varDec)
	}
	return vars, nil
}

func (p *parser) parseFunctionDeclaration() (res *ast.Function, rerr error) {
	name, err := p.identifier("function name")
	if err != nil {
		return nil, err
	}
	f := &ast.Function{Name: name}
	if p.sees('(') {
		if err := p.eat('('); err != nil {
			return nil, err
		}
		for p.sees(tokenWord) || p.sees('[') {
			param := &ast.Parameter{}

			if p.seesWord("var") {
				if err := p.eatWord("var"); err != nil {
					return nil, err
				}
				param.Qualifier = ast.Var
			} else if p.seesWord("const") {
				if err := p.eatWord("const"); err != nil {
					return nil, err
				}
				param.Qualifier = ast.Const
				if p.sees('[') {
					if err := p.eat('['); err != nil {
						return nil, err
					}
					if err := p.eatWord("ref"); err != nil {
						return nil, err
					}
					if err := p.eat(']'); err != nil {
						return nil, err
					}
					param.Qualifier = ast.ConstRef
				}
			} else if p.seesWord("out") {
				if err := p.eatWord("out"); err != nil {
					return nil, err
				}
				param.Qualifier = ast.Out
			} else if p.sees('[') {
				if err := p.eat('['); err != nil {
					return nil, err
				}
				if err := p.eatWord("ref"); err != nil {
					return nil, err
				}
				if err := p.eat(']'); err != nil {
					return nil, err
				}
				if err := p.eatWord("const"); err != nil {
					return nil, err
				}
				param.Qualifier = ast.RefConst
			}

			firstId, err := p.identifier("parameter name")
			if err != nil {
				return nil, err
			}
			param.Names = append(param.Names, firstId)
			for p.sees(',') {
				if err := p.eat(','); err != nil {
					return nil, err
				}
				id, err := p.identifier("parameter name")
				if err != nil {
					return nil, err
				}
				param.Names = append(param.Names, id)
			}
			if p.sees(':') {
				if err := p.eat(':'); err != nil {
					return nil, err
				}
				pt, err := p.qualifiedIdentifier("parameter type")
				if err != nil {
					return nil, err
				}
				param.Type = pt
			}
			f.Parameters = append(f.Parameters, param)
			if p.sees(';') {
				p.eat(';')
			} else if p.sees(')') {
				break
			} else {
				if err := p.eat(','); err != nil {
					return nil, err
				}
				break // The last parameter is not followed by a ';'.
			}
		}
		if err := p.eat(')'); err != nil {
			return nil, err
		}
	}
	if p.sees(':') {
		if err := p.eat(':'); err != nil {
			return nil, err
		}
		rt, err := p.qualifiedIdentifier("return type")
		if err != nil {
			return nil, err
		}
		f.Returns = rt
	}
	if err := p.eat(';'); err != nil {
		return nil, err
	}
	return f, nil
}

func (p *parser) parseVariableDeclaration() (*ast.Variable, error) {
	name, err := p.identifier("field name")
	if err != nil {
		return nil, err
	}
	if err := p.eat(':'); err != nil {
		return nil, err
	}
	typ, err := p.qualifiedIdentifier("type name")
	if err != nil {
		return nil, err
	}
	if err := p.eat(';'); err != nil {
		return nil, err
	}
	return &ast.Variable{Name: name, Type: typ}, nil
}

func (p *parser) nextToken() token {
	if p.isPeeking {
		// Remove the queued token from our peek queue.
		p.isPeeking = false
		return p.peekingAt
	}

	// Find the next token which is not a white-space.
	t := p.tokens.next()
	for t.tokenType == tokenWhiteSpace || t.tokenType == tokenComment {
		t = p.tokens.next()
	}
	return t
}

func (p *parser) peekToken() token {
	if !p.isPeeking {
		p.peekingAt = p.nextToken()
		p.isPeeking = true
	}
	return p.peekingAt
}

func (p *parser) sees(typ tokenType) bool {
	t := p.peekToken()
	return t.tokenType == typ
}

func (p *parser) seesWord(text string) bool {
	t := p.peekToken()
	return t.tokenType == tokenWord && strings.ToLower(t.text) == text
}

func (p *parser) seesWords(texts ...string) bool {
	for _, t := range texts {
		if p.seesWord(t) {
			return true
		}
	}
	return false
}

func (p *parser) seesKeyword() bool {
	t := p.peekToken()
	return t.tokenType == tokenWord && isKeyword(strings.ToLower(t.text))
}

func isKeyword(s string) bool {
	// TODO Complete the list of keywords, these end blocks (var, type, ...).
	return s == "implementation" || s == "var"
}

func (p *parser) eat(typ tokenType) error {
	t := p.nextToken()
	if t.tokenType != typ {
		return p.tokenError(t, typ.String())
	}
	return nil
}

func (p *parser) eats(types ...tokenType) error {
	for _, typ := range types {
		if err := p.eat(typ); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) eatWord(text string) error {
	t := p.nextToken()
	if !(t.tokenType == tokenWord && strings.ToLower(t.text) == text) {
		return p.tokenError(t, `keyword "`+text+`"`)
	}
	return nil
}

// qualifiedIdentifier parses identifiers with dots in them, e.g.
//
//     Systems.Generics.Collections
//
// There might be comments or white space between the identifiers and dots.
func (p *parser) qualifiedIdentifier(description string) (string, error) {
	s, err := p.identifier(description)
	if err != nil {
		return "", err
	}
	dot := tokenType('.')
	for p.sees(dot) {
		if err := p.eat(dot); err != nil {
			return "", err
		}
		id, err := p.identifier(description)
		if err != nil {
			return "", err
		}
		s += "." + id
	}
	return s, nil
}

func (p *parser) identifier(description string) (string, error) {
	t := p.nextToken()
	if t.tokenType == tokenWord {
		return t.text, nil
	}
	return "", p.tokenError(t, description)
}

func (p *parser) tokenError(t token, expected string) error {
	return errors.Errorf("%s expected but was %s", expected, t.String())
}
