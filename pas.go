package pas

func ParseString(code string) (*File, error) {
	return newParser([]rune(code)).parseFile()
}

type File struct {
	Kind FileKind
	Name string
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
