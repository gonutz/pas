package ast

func (*Class) isTypeDeclaration() {}

type Class struct {
	Name         string
	SuperClasses []string
	Sections     []ClassSection
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

type ClassMember interface {
	isClassMember()
}
