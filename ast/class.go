package ast

// https://docwiki.embarcadero.com/RADStudio/Alexandria/ja/%E3%82%AF%E3%83%A9%E3%82%B9%E3%81%A8%E3%82%AA%E3%83%96%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%EF%BC%88Delphi%EF%BC%89

func (*Class) isTypeDeclaration() {}

type Class struct {
	Name string
	ClassExpr
}

func NewClass(name string, superClasses ...string) *Class {
	return &Class{
		Name: name,
		ClassExpr: ClassExpr{
			SuperClasses: superClasses,
		},
	}
}
func (c *Class) WithSection(sections ...ClassSection) *Class {
	c.Sections = sections
	return c
}

type ClassExpr struct {
	SuperClasses []string
	Sections     []ClassSection
	Abstract     bool
	Sealed       bool
}

func (c *Class) AppendMemberToCurrentSection(member ClassMember) {
	if len(c.Sections) == 0 {
		c.NewSection(DefaultPublished)
	}
	i := len(c.Sections) - 1
	c.Sections[i].Members = append(c.Sections[i].Members, member)
}

func (c *Class) NewSection(v Visibility) {
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

var visibilityToString = map[Visibility]string{
	DefaultPublished: "",
	Published:        "published",
	Public:           "public",
	Protected:        "protected",
	Private:          "private",
}

func (v Visibility) String() string {
	return visibilityToString[v]
}

type ClassMember interface {
	isClassMember()
}

func (*Field) isClassMember() {}

type Field struct {
	Variable
	Class bool
}

type MethodType int

const (
	NormalMethod MethodType = iota + 1
	Constructor
	Destructor
)

func (*Method) isClassMember() {}

type Method struct {
	Function
	Type        MethodType
	Class       bool
	Strict      bool
	Virtual     bool
	Dynamic     bool
	Override    bool
	Overload    bool
	Reintroduce bool
	Final       bool
}

func (*Property) isClassMember() {}

type Property struct {
	Variable
	Class     bool
	Indexes   Parameters
	Index     int
	Reader    string
	Writer    string
	Stored    string
	Default   string
	Nodefault bool
}
