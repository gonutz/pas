package parser

import "github.com/akm/pas/ast"

func recordProcessor(record *ast.Record) func(p *parser) error {
	return func(p *parser) error {
		if err := p.eatWord("record"); err != nil {
			return err
		}
		for !p.seesWord("end") {
			if p.seesWord("procedure") {
				if err := p.eatWord("procedure"); err != nil {
					return err
				}
				f, err := p.parseFunctionDeclaration()
				if err != nil {
					return err
				}
				record.AppendMember(f)
			} else if p.seesWord("function") {
				if err := p.eatWord("function"); err != nil {
					return err
				}
				f, err := p.parseFunctionDeclaration()
				if err != nil {
					return err
				}
				record.AppendMember(f)
			} else {
				v, err := p.parseVariableDeclaration()
				if err != nil {
					return err
				}
				record.AppendMember(v)
			}
		}
		if err := p.eatWord("end"); err != nil {
			return err
		}
		if err := p.eat(';'); err != nil {
			return err
		}
		return nil
	}
}
