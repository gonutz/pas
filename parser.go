package pas

import (
	"errors"
	"strings"
)

func newParser(code []rune) *parser {
	return &parser{tokens: newTokenizer(code)}
}

type parser struct {
	tokens tokenizer
	// isPeeking and peekingAt are a one-element queue of tokens to come. The
	// tokenizer only gives us the next token, it cannot peek so we buffer one
	// token here. See parser.nextToken and parser.peekToken.
	isPeeking bool
	peekingAt token
	file      File
}

func (p *parser) parseFile() (*File, error) {
	// For now only parse units until we have tests for other kinds.
	if err := p.eatWord("unit"); err != nil {
		return nil, err
	}
	p.file.Kind = Unit
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
	if err := p.parseFileSection(InterfaceSection); err != nil {
		return nil, err
	}

	if err := p.eatWord("implementation"); err != nil {
		return nil, err
	}
	if err := p.parseFileSection(ImplementationSection); err != nil {
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

func (p *parser) parseFileSection(kind FileSectionKind) error {
	uses, err := p.parseUses()
	if err != nil {
		return err
	}
	blocks, err := p.parseSectionBlocks()
	if err != nil {
		return err
	}
	p.file.Sections = append(p.file.Sections, FileSection{
		Kind:   kind,
		Uses:   uses,
		Blocks: blocks,
	})
	return nil
}

func (p *parser) parseUses() ([]string, error) {
	var uses []string
	if p.seesWordAndEat("uses") {
		unitName, err := p.qualifiedIdentifier("uses clause")
		if err != nil {
			return nil, err
		}
		uses = append(uses, unitName)
		for p.seesAndEat(',') {
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

func (p *parser) parseSectionBlocks() ([]FileSectionBlock, error) {
	var blocks []FileSectionBlock
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

// TODO change to return *TypeBlock
func (p *parser) parseTypeBlock() (FileSectionBlock, error) {
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
		var class Class
		class.Name = identifier
		if err := p.eatWord("class"); err != nil {
			return nil, err
		}
		if p.seesAndEat('(') {
			className, err := p.qualifiedIdentifier("parent class name")
			if err != nil {
				return nil, err
			}
			class.SuperClasses = append(class.SuperClasses, className)
			for p.seesAndEat(',') {
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
		for !p.seesWord("end") {
			if p.seesWordAndEat("published") {
				class.newSection(Published)
			} else if p.seesWordAndEat("public") {
				class.newSection(Public)
			} else if p.seesWordAndEat("protected") {
				class.newSection(Protected)
			} else if p.seesWordAndEat("private") {
				class.newSection(Private)
			} else if p.seesWordAndEat("procedure") || p.seesWordAndEat("function") {
				f, err := p.parseFunctionDeclaration()
				if err != nil {
					return nil, err
				}
				class.appendMemberToCurrentSection(f)
			} else {
				v, err := p.parseVariableDeclaration()
				if err != nil {
					return nil, err
				}
				class.appendMemberToCurrentSection(v)
			}
		}
		if err := p.eatWord("end"); err != nil {
			return nil, err
		}
		if err := p.eat(';'); err != nil {
			return nil, err
		}
		return TypeBlock{class}, nil
	} else {
		var record Record
		record.Name = identifier
		if err := p.eatWord("record"); err != nil {
			return nil, err
		}
		for !p.seesWord("end") {
			if p.seesWordAndEat("procedure") || p.seesWordAndEat("function") {
				f, err := p.parseFunctionDeclaration()
				if err != nil {
					return nil, err
				}
				record.appendMember(f)
			} else {
				v, err := p.parseVariableDeclaration()
				if err != nil {
					return nil, err
				}
				record.appendMember(v)
			}
		}
		if err := p.eatWord("end"); err != nil {
			return nil, err
		}
		if err := p.eat(';'); err != nil {
			return nil, err
		}
		return TypeBlock{record}, nil
	}
}

// TODO change to return *VarBlock
func (p *parser) parseVarBlock() (FileSectionBlock, error) {
	if err := p.eatWord("var"); err != nil {
		return nil, err
	}
	var vars VarBlock
	for p.sees(tokenWord) && !p.seesKeyword() {
		varDec, err := p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}
		vars = append(vars, varDec)
	}
	return vars, nil
}

// TODO change to return *Function
func (p *parser) parseFunctionDeclaration() (ClassMember, error) {
	var f Function
	var err error
	f.Name, err = p.identifier("function name")
	if err != nil {
		return nil, err
	}
	if p.seesAndEat('(') {
		for p.sees(tokenWord) || p.sees('[') {
			var param Parameter

			if p.seesWordAndEat("var") {
				param.Qualifier = Var
			} else if p.seesWordAndEat("const") {
				param.Qualifier = Const
				if p.seesAndEat('[') {
					if err := p.eatWord("ref"); err != nil {
						return nil, err
					}
					if err := p.eat(']'); err != nil {
						return nil, err
					}
					param.Qualifier = ConstRef
				}
			} else if p.seesWordAndEat("out") {
				param.Qualifier = Out
			} else if p.seesAndEat('[') {
				if err := p.eatWord("ref"); err != nil {
					return nil, err
				}
				if err := p.eat(']'); err != nil {
					return nil, err
				}
				if err := p.eatWord("const"); err != nil {
					return nil, err
				}
				param.Qualifier = RefConst
			}

			firstId, err := p.identifier("parameter name")
			if err != nil {
				return nil, err
			}
			param.Names = append(param.Names, firstId)
			for p.seesAndEat(',') {
				id, err := p.identifier("parameter name")
				if err != nil {
					return nil, err
				}
				param.Names = append(param.Names, id)
			}
			if p.seesAndEat(':') {
				pt, err := p.qualifiedIdentifier("parameter type")
				if err != nil {
					return nil, err
				}
				param.Type = pt
			}
			f.Parameters = append(f.Parameters, param)
			if !p.seesAndEat(';') {
				break // The last parameter is not followed by a ';'.
			}
		}
		if err := p.eat(')'); err != nil {
			return nil, err
		}
	}
	if p.seesAndEat(':') {
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

// TODO change to return *Variable
func (p *parser) parseVariableDeclaration() (Variable, error) {
	var v Variable
	var err error
	v.Name, err = p.identifier("field name")
	if err != nil {
		return v, err
	}
	if err := p.eat(':'); err != nil {
		return v, err
	}
	v.Type, err = p.qualifiedIdentifier("type name")
	if err != nil {
		return v, err
	}
	if err := p.eat(';'); err != nil {
		return v, err
	}
	return v, nil
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

func (p *parser) seesAndEat(typ tokenType) bool {
	t := p.peekToken()
	if t.tokenType == typ {
		p.nextToken()
		return true
	}
	return false
}

func (p *parser) seesWord(text string) bool {
	t := p.peekToken()
	return t.tokenType == tokenWord && strings.ToLower(t.text) == text
}

func (p *parser) seesWordAndEat(text string) bool {
	t := p.peekToken()
	if t.tokenType == tokenWord && strings.ToLower(t.text) == text {
		p.nextToken()
		return true
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
	for p.seesAndEat('.') {
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
	return errors.New(expected + " expected but was " + t.String())
}
