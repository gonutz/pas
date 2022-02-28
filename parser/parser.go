package parser

import (
	"strconv"
	"strings"

	"github.com/akm/delparser/ast"
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

func (p *parser) startEndToken(start, end tokenType, fn func() error) error {
	if err := p.eat(start); err != nil {
		return err
	}
	if err := fn(); err != nil {
		return err
	}
	if err := p.eat(end); err != nil {
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
			varBlock, err := p.parseVarBlock("var")
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, varBlock)
		} else if p.seesWord("threadvar") {
			varBlock, err := p.parseVarBlock("threadvar")
			if err != nil {
				return nil, err
			}
			threadVarBlock := make(ast.ThreadVarBlock, len(varBlock))
			for i, v := range varBlock {
				threadVarBlock[i] = v
			}
			blocks = append(blocks, threadVarBlock)
		} else if p.seesWord("function") {
			if err := p.eatWord("function"); err != nil {
				return nil, err
			}
			f, err := p.parseFunctionDeclaration()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, f)
		} else if p.seesWord("procedure") {
			if err := p.eatWord("procedure"); err != nil {
				return nil, err
			}
			f, err := p.parseFunctionDeclaration()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, f)
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
		expr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		res = append(res, &ast.Type{Name: identifier, Expr: expr})
		if p.sees(tokenWord) && p.seesReservedWord() {
			break
		}
	}
	return res, nil
}

func (p *parser) parseTypeExpr() (ast.TypeExpr, error) {
	if p.seesWord("class") {
		expr := &ast.ClassExpr{}
		if err := classProc(expr)(p); err != nil {
			return nil, err
		}
		return expr, nil
	} else if p.seesWord("record") {
		expr := &ast.RecordExpr{}
		if err := recordProc(expr)(p); err != nil {
			return nil, err
		}
		return expr, nil
	} else if p.seesWords("packed", "array") {
		expr := &ast.ArrayExpr{}
		if err := arrayProc(expr)(p); err != nil {
			return nil, err
		}
		if err := p.eat(';'); err != nil {
			return nil, err
		}
		return expr, nil
	} else {
		return nil, errors.Errorf("expected type expression, got %+v", p.peekToken())
	}
}

func (p *parser) parseVarBlock(word string) (ast.VarBlock, error) {
	if err := p.eatWord(word); err != nil {
		return nil, err
	}
	var vars ast.VarBlock
	for p.sees(tokenWord) && !p.seesReservedWord() {
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
	if err := functionProc(f)(p); err != nil {
		return nil, err
	}
	return f, nil
}

func (p *parser) parseParameters(endToken tokenType) (ast.Parameters, error) {
	res := ast.Parameters{}
	if err := parametersProc(&res, endToken)(p); err != nil {
		return nil, err
	}
	return res, nil
}

func (p *parser) parseVariableDeclaration() (*ast.Variable, error) {
	names, err := p.parseSeparatedString(',', "variable name")
	if err != nil {
		return nil, err
	}
	r := &ast.Variable{Names: names}
	if err := p.eat(':'); err != nil {
		return nil, err
	}
	if p.seesWords("packed", "array") {
		arrayExpr := &ast.ArrayExpr{}
		if err := arrayProc(arrayExpr)(p); err != nil {
			return nil, err
		}
		r.Type = arrayExpr
	} else {
		typ, err := p.qualifiedIdentifier("type name")
		if err != nil {
			return nil, err
		}
		if p.sees('[') {
			err := p.startEndToken('[', ']', func() error {
				token := p.nextToken()
				if token.tokenType != tokenInt {
					return errors.Errorf("expected int, got %+v", token)
				}
				l, err := strconv.Atoi(token.text)
				if err != nil {
					return err
				}
				r.Length = l
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
		r.Type = ast.TypeName(typ)
	}
	if p.seesWord("absolute") {
		if err := p.eatWord("absolute"); err != nil {
			return nil, err
		}
		abs, err := p.identifier("absolute reference name")
		if err != nil {
			return nil, err
		}
		r.Absolute = abs
	}

	if p.sees('=') {
		if err := p.eat('='); err != nil {
			return nil, err
		}
		t := p.nextToken()
		switch t.tokenType {
		case tokenWord, tokenInt, tokenReal, tokenString:
			// OK
		default:
			return nil, errors.Errorf("expected parameter default value, got %+v", t)
		}
		r.Default = t.text
	}
	if err := p.eat(';'); err != nil {
		return nil, err
	}
	return r, nil
}

func (p *parser) parseProperty() (*ast.Property, error) {
	name, err := p.identifier("property name")
	if err != nil {
		return nil, err
	}
	res := &ast.Property{Variable: ast.Variable{Names: []string{name}}}
	if err := propertyProc(res)(p); err != nil {
		return nil, err
	}
	return res, nil
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

func (p *parser) seesReservedWord() bool {
	t := p.peekToken()
	return t.tokenType == tokenWord && isReservedWord(strings.ToLower(t.text))
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
