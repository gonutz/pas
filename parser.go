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
	err    error
}

func (p *parser) parseFile() (*File, error) {
	var unit File // For now only parse units until we have tests for other kinds.
	p.eatWord("unit")
	unit.Kind = Unit
	unit.Name = p.identifier("unit name")
	p.eat(';')
	p.eatWord("interface")
	p.eatWord("implementation")
	p.eatWord("end")
	p.eat('.')
	return &unit, p.err
}

func (p *parser) nextToken() token {
	t := p.tokens.next()
	for t.tokenType == tokenWhiteSpace {
		t = p.tokens.next()
	}
	return t
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
		p.tokenError(t, "keyword "+text)
	}
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
	p.err = errors.New(expected + " expected but was " + t.tokenType.String())
}
