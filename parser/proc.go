package parser

type namedProc struct {
	name string
	fn   func(*parser) error
}

type procSelector struct {
	procs       []*namedProc
	defaultProc func(*parser) error
}

func (fs *procSelector) Do(p *parser) error {
	processed := false
	for _, proc := range fs.procs {
		if p.seesWord(proc.name) {
			if err := p.eatWord(proc.name); err != nil {
				return err
			}
			if err := proc.fn(p); err != nil {
				return err
			}
			processed = true
			break
		}
	}
	if !processed {
		if err := fs.defaultProc(p); err != nil {
			return err
		}
	}
	return nil
}
