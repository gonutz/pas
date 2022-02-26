package parser

import "github.com/akm/pas/ast"

func classProcessor(class *ast.Class) func(p *parser) error {
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
		class.AppendMemberToCurrentSection(&ast.Method{
			Function: *f,
		})
		return nil
	}
	selector := &procSelector{
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
			class.AppendMemberToCurrentSection(&ast.Field{Variable: *v})
			return nil
		},
	}
	return func(p *parser) error {
		if err := p.eatWord("class"); err != nil {
			return err
		}
		if p.sees('(') {
			if err := p.eat('('); err != nil {
				return err
			}
			superClassNames, err := p.parseSeparatedString(',', "parent class name", "parent interface name")
			if err != nil {
				return err
			}
			class.SuperClasses = superClassNames
			if err := p.eat(')'); err != nil {
				return err
			}
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
