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
	uses := p.parseUses()
	blocks := p.parseSectionBlocks()
	p.file.Sections = append(p.file.Sections, FileSection{
		Kind:   kind,
		Uses:   uses,
		Blocks: blocks,
	})
}

func (p *parser) parseUses() []string {
	var uses []string
	if p.seesWordAndEat("uses") {
		uses = append(uses, p.qualifiedIdentifier("uses clause"))
		for p.seesAndEat(',') {
			uses = append(uses, p.qualifiedIdentifier("uses clause"))
		}
		p.eat(';')
	}
	return uses
}

func (p *parser) parseSectionBlocks() []FileSectionBlock {
	var blocks []FileSectionBlock
	for {
		if p.seesWord("type") {
			blocks = append(blocks, p.parseTypeBlock())
		} else if p.seesWord("var") {
			blocks = append(blocks, p.parseVarBlock())
		} else {
			break
		}
	}
	return blocks
}

func (p *parser) parseTypeBlock() FileSectionBlock {
	p.eatWord("type")
	identifier := p.identifier("type name")
	p.eat('=')
	if p.seesWord("class") {
		var class Class
		class.Name = identifier
		p.eatWord("class")
		if p.seesAndEat('(') {
			class.SuperClasses = append(
				class.SuperClasses,
				p.qualifiedIdentifier("parent class name"),
			)
			for p.seesAndEat(',') {
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
		return TypeBlock{class}
	} else {
		var record Record
		record.Name = identifier
		p.eatWord("record")
		for !(p.seesWord("end") || p.err != nil) {
			if p.seesWordAndEat("procedure") || p.seesWordAndEat("function") {
				record.appendMember(p.parseFunctionDeclaration())
			} else {
				record.appendMember(p.parseVariableDeclaration())
			}
		}
		p.eatWord("end")
		p.eat(';')
		return TypeBlock{record}
	}
}

func (p *parser) parseVarBlock() FileSectionBlock {
	p.eatWord("var")
	var vars VarBlock
	for p.sees(tokenWord) && !p.seesKeyword() {
		vars = append(vars, p.parseVariableDeclaration())
	}
	return vars
}

func (p *parser) parseFunctionDeclaration() ClassMember {
	var f Function
	f.Name = p.identifier("function name")
	if p.seesAndEat('(') {
		for p.sees(tokenWord) || p.sees('[') {
			var param Parameter

			if p.seesWordAndEat("var") {
				param.Qualifier = Var
			} else if p.seesWordAndEat("const") {
				param.Qualifier = Const
				if p.seesAndEat('[') {
					p.eatWord("ref")
					p.eat(']')
					param.Qualifier = ConstRef
				}
			} else if p.seesWordAndEat("out") {
				param.Qualifier = Out
			} else if p.seesAndEat('[') {
				p.eatWord("ref")
				p.eat(']')
				p.eatWord("const")
				param.Qualifier = RefConst
			}

			param.Names = append(param.Names, p.identifier("parameter name"))
			for p.seesAndEat(',') {
				param.Names = append(param.Names, p.identifier("parameter name"))
			}
			if p.seesAndEat(':') {
				param.Type = p.qualifiedIdentifier("parameter type")
			}
			f.Parameters = append(f.Parameters, param)
			if !p.seesAndEat(';') {
				break // The last parameter is not followed by a ';'.
			}
		}
		p.eat(')')
	}
	if p.seesAndEat(':') {
		f.Returns = p.qualifiedIdentifier("return type")
	}
	p.eat(';')
	return f
}

func (p *parser) parseVariableDeclaration() Variable {
	var v Variable
	v.Name = p.identifier("field name")
	p.eat(':')
	v.Type = p.qualifiedIdentifier("type name")
	p.eat(';')
	return v
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
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	return t.tokenType == typ
}

func (p *parser) seesAndEat(typ tokenType) bool {
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	if t.tokenType == typ {
		p.nextToken()
		return true
	}
	return false
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

func (p *parser) seesKeyword() bool {
	if p.err != nil {
		return false
	}
	t := p.peekToken()
	return t.tokenType == tokenWord && isKeyword(strings.ToLower(t.text))
}

func isKeyword(s string) bool {
	// TODO Complete the list of keywords, these end blocks (var, type, ...).
	return s == "implementation" || s == "var"
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
	for p.seesAndEat('.') {
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
