package ast

func (VarBlock) isFileSectionBlock() {}

type VarBlock []*Variable

func (*Variable) isClassMember() {}

type Variable struct {
	Name string
	Type string
}
