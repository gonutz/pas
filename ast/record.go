package ast

func (*Record) isTypeDeclaration() {}

type Record struct {
	Name string
	RecordExpr
}

func NewRecord(name string, members ...RecordMember) *Record {
	return &Record{
		Name: name,
		RecordExpr: RecordExpr{
			Members: members,
		},
	}
}

func (*RecordExpr) isVarType()  {}
func (*RecordExpr) isTypeExpr() {}

type RecordExpr struct {
	Members []RecordMember
}

// See https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E6%A7%8B%E9%80%A0%E5%8C%96%E5%9E%8B%EF%BC%88Delphi%EF%BC%89#.E3.83.AC.E3.82.B3.E3.83.BC.E3.83.89.E5.9E.8B.EF.BC.88.E9.AB.98.E5.BA.A6.EF.BC.89
type RecordMember = ClassMember

func (c *Record) AppendMember(member RecordMember) {
	c.Members = append(c.Members, member)
}
