package ast

func (VarBlock) isFileSectionBlock() {}

type VarBlock []*Variable

func (ThreadVarBlock) isFileSectionBlock() {}

type ThreadVarBlock []*Variable

type Variable struct {
	Names    []string
	Type     VarType
	Default  string
	Length   int
	Absolute string
}

func NewVariable(name string, nameOrType string, nameOrTypes ...string) *Variable {
	var names []string
	var typ string
	if len(nameOrTypes) == 0 {
		names = []string{name}
		typ = nameOrType
	} else {
		l := len(nameOrTypes)
		typ = nameOrTypes[l-1]
		names = append([]string{name, nameOrType}, nameOrTypes[:l-1]...)
	}
	return &Variable{Names: names, Type: TypeName(typ)}
}

func (v *Variable) WithDefault(s string) *Variable {
	v.Default = s
	return v
}
func (v *Variable) WithLength(val int) *Variable {
	v.Length = val
	return v
}
func (v *Variable) WithAbsolute(s string) *Variable {
	v.Absolute = s
	return v
}

type VarType interface {
	isVarType()
}

func (TypeName) isVarType() {}

type TypeName string
