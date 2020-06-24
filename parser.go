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
	err       error
}

func (p *parser) parseFile() (*File, error) {
	// For now only parse units until we have tests for other kinds.
	p.eatWord("unit")
	p.file.Kind = Unit
	p.file.Name = p.qualifiedIdentifier("unit name")
	p.eat(';')

	p.eatWord("interface")
	p.parseFileSection(InterfaceSection)

	p.eatWord("implementation")
	p.parseFileSection(ImplementationSection)

	p.eatWord("end")
	p.eat('.')
	return &p.file, p.err
}

func (p *parser) parseFileSection(kind FileSectionKind) {
	p.file.Sections = append(p.file.Sections, FileSection{
		Kind:   kind,
		Uses:   p.parseOptionalUses(),
		Blocks: p.parseSectionBlocks(),
	})
}

func (p *parser) parseSectionBlocks() []FileSectionBlock {
	var blocks []FileSectionBlock
	if p.seesWordAndEat("type") {
		var class Class
		class.Name = p.identifier("type name")
		p.eat('=')
		p.eatWord("class")
		if p.sees('(') {
			p.eat('(')
			class.SuperClasses = append(
				class.SuperClasses,
				p.qualifiedIdentifier("parent class name"),
			)
			for p.sees(',') {
				p.eat(',')
				class.SuperClasses = append(
					class.SuperClasses,
					p.qualifiedIdentifier("parent interface name"),
				)
			}
			p.eat(')')
		}
		for !(p.seesWord("end") || p.err != nil) {
			if p.seesWordAndEat("published") {
				class.newSection(Published)
			} else if p.seesWordAndEat("public") {
				class.newSection(Public)
			} else if p.seesWordAndEat("protected") {
				class.newSection(Protected)
			} else if p.seesWordAndEat("private") {
				class.newSection(Private)
			} else if p.seesWordAndEat("procedure") || p.seesWordAndEat("function") {
				class.appendMemberToCurrentSection(p.parseFunctionDeclaration())
			} else {
				class.appendMemberToCurrentSection(p.parseVariableDeclaration())
			}
		}
		p.eatWord("end")
		p.eat(';')
		blocks = append(blocks, TypeBlock{class})
	}
	return blocks
}

func (p *parser) parseFunctionDeclaration() ClassMember {
	var f Function
	f.Name = p.identifier("function name")
	if p.sees('(') {
		p.eat('(')
		for p.sees(tokenWord) || p.sees('[') {
			var param Parameter

			if p.seesWordAndEat("var") {
				param.Qualifier = Var
			} else if p.seesWordAndEat("const") {
				param.Qualifier = Const
				if p.sees('[') {
					p.eat('[')
					p.eatWord("ref")
					p.eat(']')
					param.Qualifier = ConstRef
				}
			} else if p.seesWordAndEat("out") {
				param.Qualifier = Out
			} else if p.sees('[') {
				p.eat('[')
				p.eatWord("ref")
				p.eat(']')
				p.eatWord("const")
				param.Qualifier = RefConst
			}

			param.Names = append(param.Names, p.identifier("parameter name"))
			for p.sees(',') {
				p.eat(',')
				param.Names = append(param.Names, p.identifier("parameter name"))
			}
			if p.sees(':') {
				p.eat(':')
				param.Type = p.qualifiedIdentifier("parameter type")
			}
			f.Parameters = append(f.Parameters, param)
			if p.sees(';') {
				p.eat(';')
			} else {
				break
			}
		}
		p.eat(')')
	}
	if p.sees(':') {
		p.eat(':')
		f.Returns = p.qualifiedIdentifier("return type")
	}
	p.eat(';')
	return f
}

func (p *parser) parseVariableDeclaration() ClassMember {
	var v Variable
	v.Name = p.identifier("field name")
	p.eat(':')
	v.Type = p.qualifiedIdentifier("type name")
	p.eat(';')
	return v
}

func (p *parser) parseOptionalUses() []string {
	var uses []string
	if p.seesWordAndEat("uses") {
		uses = append(uses, p.qualifiedIdentifier("uses clause"))
		for p.sees(',') {
			p.eat(',')
			uses = append(uses, p.qualifiedIdentifier("uses clause"))
		}
		p.eat(';')
	}
	return uses
}

func (p *parser) nextToken() token {
	if p.isPeeking {
		// Remove the queued token from our peek queue.
		p.isPeeking = false
		return p.peekingAt
	}

	// Find the next token which is not a white-space.
	t := p.tokens.next()
	for t.tokenType == tokenWhiteSpace {
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
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	return t.tokenType == typ
}

func (p *parser) seesWord(text string) bool {
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	return t.tokenType == tokenWord && strings.ToLower(t.text) == text
}

func (p *parser) seesWordAndEat(text string) bool {
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	if t.tokenType == tokenWord && strings.ToLower(t.text) == text {
		p.nextToken()
		return true
	}
	return false
}

func (p *parser) eat(typ tokenType) {
	if p.err != nil {
		return
	}
	t := p.nextToken()
	if t.tokenType != typ {
		p.tokenError(t, typ.String())
	}
}

func (p *parser) eatWord(text string) {
	if p.err != nil {
		return
	}
	t := p.nextToken()
	if !(t.tokenType == tokenWord && strings.ToLower(t.text) == text) {
		p.tokenError(t, `keyword "`+text+`"`)
	}
}

// qualifiedIdentifier parses identifiers with dots in them, e.g.
//
//     Systems.Generics.Collections
//
// There might be comments or white space between the identifiers and dots.
func (p *parser) qualifiedIdentifier(description string) string {
	s := p.identifier(description)
	for p.sees('.') {
		p.eat('.')
		s += "." + p.identifier(description)
	}
	return s
}

func (p *parser) identifier(description string) string {
	if p.err != nil {
		return ""
	}
	t := p.nextToken()
	if t.tokenType == tokenWord {
		return t.text
	}
	p.tokenError(t, description)
	return ""
}

func (p *parser) tokenError(t token, expected string) {
	p.err = errors.New(expected + " expected but was " + t.String())
}
