package ast

func (*Function) isClassMember() {}

type Function struct {
	Name       string
	Parameters []*Parameter
	// Returns is either the return type for functions or the empty string for
	// procedures.
	Returns string
}

type Parameter struct {
	Names []string
	// Type might be empty. In that case this is an untyped parameter like in:
	//
	//     procedure(const A; var B);
	Type      string
	Qualifier Qualifier
}

type Qualifier int

const (
	NoQualifier           = 0
	Var         Qualifier = 1
	Const       Qualifier = 2
	// ConstRef is "const [Ref]" which semantically is the same as RefConst.
	ConstRef Qualifier = 3
	// RefConst is "[Ref] const" which semantically is the same as ConstRef.
	RefConst Qualifier = 4
	Out      Qualifier = 5
)

func (q Qualifier) String() string {
	switch q {
	case NoQualifier:
		return ""
	case Var:
		return "var"
	case Const:
		return "const"
	case ConstRef:
		return "const [Ref]"
	case RefConst:
		return "[Ref] const"
	case Out:
		return "out"
	}
	return "unknown Qualifier"
}
