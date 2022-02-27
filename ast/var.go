package ast

func (VarBlock) isFileSectionBlock() {}

type VarBlock []*Variable

type Variable struct {
	Name string
	Type string
}
