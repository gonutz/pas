package ast

func (*EnumExpr) isTypeExpr() {}
func (*EnumExpr) isVarType()  {}

// https://docwiki.embarcadero.com/RADStudio/Sydney/en/Simple_Types_(Delphi)#Enumerated_Types
type EnumExpr struct {
	Members []EnumMember
}

type EnumMember struct {
	Name  string
	Value string
}
