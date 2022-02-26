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

func (p *parser) parseUses() ([]string, error) {
	var uses []string
	if p.seesWord("uses") {
		if err := p.eatWord("uses"); err != nil {
			return nil, err
		}
		unitName, err := p.qualifiedIdentifier("uses clause")
		if err != nil {
			return nil, err
		}
		uses = append(uses, unitName)
		for p.sees(',') {
			if err := p.eat(','); err != nil {
				return nil, err
			}
			unitName, err := p.qualifiedIdentifier("uses clause")
			if err != nil {
				return nil, err
			}
			uses = append(uses, unitName)
		}
		if err := p.eat(';'); err != nil {
			return nil, err
		}
	}
	return uses, nil
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
	identifier, err := p.identifier("type name")
	if err != nil {
		return nil, err
	}
	if err := p.eat('='); err != nil {
		return nil, err
	}
	if p.seesWord("class") {
		class, err := p.parseClass(identifier)
		if err != nil {
			return nil, err
		}
		return ast.TypeBlock{class}, nil
	} else {
		record, err := p.parseRecord(identifier)
		if err != nil {
			return nil, err
		}
		return ast.TypeBlock{record}, nil
	}
}

type namedProc struct {
	name string
	fn   func(*parser, string) error
}

type procSelector struct {
	procs       []*namedProc
	defaultProc func(*parser) error
}

func (fs *procSelector) Do(p *parser) error {
	processed := false
	for _, proc := range fs.procs {
		if p.seesWord(proc.name) {
			if err := proc.fn(p, proc.name); err != nil {
				return err
			}
			processed = true
			break
		}
	}
	if !processed {
		if err := fs.defaultProc(p); err != nil {
			return err
		}
	}
	return nil
}

func classMemberProcessor(class *ast.Class) *procSelector {
	newSection := func(visibility ast.Visibility) func(p *parser, name string) error {
		return func(p *parser, name string) error {
			if err := p.eatWord(name); err != nil {
				return err
			}
			class.NewSection(visibility)
			return nil
		}
	}
	appendFunc := func(p *parser, name string) error {
		if err := p.eatWord(name); err != nil {
			return err
		}
		f, err := p.parseFunctionDeclaration()
		if err != nil {
			return err
		}
		class.AppendMemberToCurrentSection(f)
		return nil
	}
	return &procSelector{
		procs: []*namedProc{
			{"published", newSection(ast.Published)},
			{"public", newSection(ast.Public)},
			{"protected", newSection(ast.Protected)},
			{"private", newSection(ast.Private)},
			{"procedure", appendFunc},
			{"function", appendFunc},
		},
		defaultProc: func(p *parser) error {
			v, err := p.parseVariableDeclaration()
			if err != nil {
				return err
			}
			class.AppendMemberToCurrentSection(v)
			return nil
		},
	}
}

func (p *parser) parseClass(identifier string) (*ast.Class, error) {
	class := &ast.Class{Name: identifier}
	if err := p.eatWord("class"); err != nil {
		return nil, err
	}
	if p.sees('(') {
		if err := p.eat('('); err != nil {
			return nil, err
		}
		className, err := p.qualifiedIdentifier("parent class name")
		if err != nil {
			return nil, err
		}
		class.SuperClasses = append(class.SuperClasses, className)
		for p.sees(',') {
			if err := p.eat(','); err != nil {
				return nil, err
			}
			intf, err := p.qualifiedIdentifier("parent interface name")
			if err != nil {
				return nil, err
			}
			class.SuperClasses = append(class.SuperClasses, intf)
		}
		if err := p.eat(')'); err != nil {
			return nil, err
		}
	}

	proc := classMemberProcessor(class)
	for !p.seesWord("end") {
		if err := proc.Do(p); err != nil {
			return nil, err
		}
	}
	if err := p.eatWord("end"); err != nil {
		return nil, err
	}
	if err := p.eat(';'); err != nil {
		return nil, err
	}
	return class, nil
}

func (p *parser) parseRecord(identifier string) (*ast.Record, error) {
	record := &ast.Record{Name: identifier}
	if err := p.eatWord("record"); err != nil {
		return nil, err
	}
	for !p.seesWord("end") {
		if p.seesWord("procedure") {
			if err := p.eatWord("procedure"); err != nil {
				return nil, err
			}
			f, err := p.parseFunctionDeclaration()
			if err != nil {
				return nil, err
			}
			record.AppendMember(f)
		} else if p.seesWord("function") {
			if err := p.eatWord("function"); err != nil {
				return nil, err
			}
			f, err := p.parseFunctionDeclaration()
			if err != nil {
				return nil, err
			}
			record.AppendMember(f)
		} else {
			v, err := p.parseVariableDeclaration()
			if err != nil {
				return nil, err
			}
			record.AppendMember(v)
		}
	}
	if err := p.eatWord("end"); err != nil {
		return nil, err
	}
	if err := p.eat(';'); err != nil {
		return nil, err
	}
	return record, nil
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
