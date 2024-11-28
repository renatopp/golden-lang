package golang

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/codegen/core"
)

type scoped struct {
	instructions []core.Instruction
}

func newScoped() *scoped {
	return &scoped{
		instructions: []core.Instruction{},
	}
}

func (s *scoped) Print(depth int) string {
	res := ""
	for _, i := range s.instructions {
		res += i.Print(depth) + "\n"
	}
	return res
}

func (s *scoped) Append(i core.Instruction) {
	s.instructions = append(s.instructions, i)
}

//
//
//

const package_go = `package %s

%s

%s
`
const package_go_import = `import "%s"`

type Package struct {
	*scoped
	PackageName string
	PackagePath string
	Imports     []string
}

func (p *Package) Print(depth int) string {
	importList := []string{}
	for _, i := range p.Imports {
		importList = append(importList, fmt.Sprintf(package_go_import, i))
	}
	imports := strings.Join(importList, "\n")

	instructionList := []string{}
	for _, i := range p.instructions {
		instructionList = append(instructionList, i.Print(depth+1))
	}
	instructions := strings.Join(instructionList, "\n")

	res := fmt.Sprintf(package_go, p.PackageName, imports, instructions)
	return res
}
