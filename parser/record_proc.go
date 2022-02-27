package parser

import "github.com/akm/pas/ast"

func recordProc(record *ast.Record) func(p *parser) error {
	appendFunc := func(p *parser) error {
		f, err := p.parseFunctionDeclaration()
		if err != nil {
			return err
		}
		record.AppendMember(&ast.Method{Function: *f})
		return nil
	}

	selector := &procSelector{
		procs: []*namedProc{
			{"procedure", appendFunc},
			{"function", appendFunc},
		},
		defaultProc: func(p *parser) error {
			v, err := p.parseVariableDeclaration()
			if err != nil {
				return err
			}
			record.AppendMember(&ast.Field{Variable: *v})
			return nil
		},
	}

	return func(p *parser) error {
		if err := p.eatWord("record"); err != nil {
			return err
		}
		for !p.seesWord("end") {
			if err := selector.Do(p); err != nil {
				return err
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
