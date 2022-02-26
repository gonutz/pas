package pas

func ParseString(code string) (*File, error) {
	return newParser([]rune(code)).parseFile()
}

type File struct {
	Kind     FileKind
	Name     string
	Sections []*FileSection
}

type FileKind int

const (
	Program FileKind = 0
	Unit    FileKind = 1
	Library FileKind = 2
	Package FileKind = 3
)

func (k FileKind) String() string {
	if k == Program {
		return "program"
	} else if k == Unit {
		return "unit"
	} else if k == Library {
		return "library"
	} else if k == Package {
		return "package"
	}
	return "unknown FileKind"
}

type FileSection struct {
	Kind   FileSectionKind
	Uses   []string
	Blocks []FileSectionBlock
}

type FileSectionKind int

const (
	InterfaceSection      FileSectionKind = 0
	ImplementationSection FileSectionKind = 1
	InitializationSection FileSectionKind = 2
	FinalizationSection   FileSectionKind = 3
)

func (k FileSectionKind) String() string {
	if k == InterfaceSection {
		return "interface"
	} else if k == ImplementationSection {
		return "implementation"
	} else if k == InitializationSection {
		return "initialization"
	} else if k == FinalizationSection {
		return "finalization"
	}
	return "unknown FileSectionKind"
}

type FileSectionBlock interface {
	isFileSectionBlock()
}

func (TypeBlock) isFileSectionBlock() {}
func (VarBlock) isFileSectionBlock()  {}

type TypeBlock []TypeDeclaration

type VarBlock []*Variable

type TypeDeclaration interface {
	isTypeDeclaration()
}

func (*Class) isTypeDeclaration() {}

type Class struct {
	Name         string
	SuperClasses []string
	Sections     []ClassSection
}

func (c *Class) appendMemberToCurrentSection(member ClassMember) {
	if len(c.Sections) == 0 {
		c.newSection(DefaultPublished)
	}
	i := len(c.Sections) - 1
	c.Sections[i].Members = append(c.Sections[i].Members, member)
}

func (c *Class) newSection(v Visibility) {
	c.Sections = append(c.Sections, ClassSection{Visibility: v})
}

type ClassSection struct {
	Visibility Visibility
	Members    []ClassMember
}

type Visibility int

const (
	// DefaultPublished is the unnamed first section in a class, in Delphi the
	// visibility of a class defaults to published if you do not specify it.
	DefaultPublished Visibility = 0
	Published        Visibility = 1
	Public           Visibility = 2
	Protected        Visibility = 3
	Private          Visibility = 4
)

type ClassMember interface {
	isClassMember()
}

func (*Variable) isClassMember() {}
func (*Function) isClassMember() {}

type Variable struct {
	Name string
	Type string
}

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

func (*Record) isTypeDeclaration() {}

type Record struct {
	Name string
	// Sections []RecordSection
	Members []RecordMember
}

// See https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E6%A7%8B%E9%80%A0%E5%8C%96%E5%9E%8B%EF%BC%88Delphi%EF%BC%89#.E3.83.AC.E3.82.B3.E3.83.BC.E3.83.89.E5.9E.8B.EF.BC.88.E9.AB.98.E5.BA.A6.EF.BC.89
type RecordMember = ClassMember

func (c *Record) appendMember(member RecordMember) {
	c.Members = append(c.Members, member)
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
