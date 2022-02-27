package parser

import "github.com/akm/pas/ast"

func classProcessor(class *ast.Class) func(p *parser) error {
	newSection := func(visibility ast.Visibility) func(p *parser) error {
		return func(p *parser) error {
			class.NewSection(visibility)
			return nil
		}
	}
	appendMethod := func(static bool) func(p *parser) error {
		return func(p *parser) error {
			f, err := p.parseFunctionDeclaration()
			if err != nil {
				return err
			}
			class.AppendMemberToCurrentSection(&ast.Method{
				Class:    static,
				Function: *f,
			})
			return nil
		}
	}

	appendField := func(static bool) func(p *parser) error {
		return func(p *parser) error {
			v, err := p.parseVariableDeclaration()
			if err != nil {
				return err
			}
			class.AppendMemberToCurrentSection(&ast.Field{
				Class:    static,
				Variable: *v,
			})
			return nil
		}
	}

	appendProperty := func(static bool) func(p *parser) error {
		return func(p *parser) error {
			prop, err := p.parseProperty()
			if err != nil {
				return err
			}
			prop.Class = static
			class.AppendMemberToCurrentSection(prop)
			return nil
		}
	}

	classSelector := &procSelector{
		procs: []*namedProc{
			{"procedure", appendMethod(true)},
			{"function", appendMethod(true)},
			{"property", appendProperty(true)},
			{"var", appendField(true)},
		},
		defaultProc: appendField(true),
	}

	selector := &procSelector{
		procs: []*namedProc{
			{"published", newSection(ast.Published)},
			{"public", newSection(ast.Public)},
			{"protected", newSection(ast.Protected)},
			{"private", newSection(ast.Private)},
			{"procedure", appendMethod(false)},
			{"function", appendMethod(false)},
			{"property", appendProperty(false)},
			{"class", classSelector.Do},
		},
		defaultProc: appendField(false),
	}
	return func(p *parser) error {
		if err := p.eatWord("class"); err != nil {
			return err
		}
		if p.sees('(') {
			err := p.startEndToken('(', ')', func() error {
				superClassNames, err := p.parseSeparatedString(',', "parent class name", "parent interface name")
				if err != nil {
					return err
				}
				class.SuperClasses = superClassNames
				return nil
			})
			if err != nil {
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
