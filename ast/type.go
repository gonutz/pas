package ast

func (TypeBlock) isFileSectionBlock() {}

type TypeBlock []*Type

type Type struct {
	Name string
	Expr TypeExpr
}

type TypeExpr interface {
	isTypeExpr()
	VarType
}
