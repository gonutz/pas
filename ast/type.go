package ast

func (TypeBlock) isFileSectionBlock() {}

type TypeBlock []*Type

type Type struct {
	Name string
	Type TypeExpr
}

type TypeExpr interface {
	isTypeExpr()
	VarType
}
