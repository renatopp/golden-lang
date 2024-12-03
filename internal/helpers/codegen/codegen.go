package codegen

import "strings"

func JoinList[T any](separator string, list []T, f func(T) string) string {
	s := ""
	for i, item := range list {
		if i > 0 {
			s += separator
		}
		s += f(item)
	}
	return s
}

//
//
//

type Identer struct {
	identLevel int
}

func NewIdenter() *Identer {
	return &Identer{
		identLevel: 0,
	}
}

func (i *Identer) Inc() { i.identLevel++ }
func (i *Identer) Dec() { i.identLevel-- }
func (i *Identer) Indent(block string) string {
	spaces := strings.Repeat(" ", i.identLevel*2)
	lineSpace := "\n" + spaces
	return spaces + strings.ReplaceAll(block, "\n", lineSpace)
}

//
//
//
