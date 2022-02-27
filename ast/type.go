package ast

func (TypeBlock) isFileSectionBlock() {}

type TypeBlock []TypeDeclaration
type TypeDeclaration interface {
	isTypeDeclaration()
}

type TypeExpr interface {
	isTypeExpr()
	VarType
}
