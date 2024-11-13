package build

import "github.com/renatopp/golden/internal/helpers/fs"

type StepDiscoverPackage struct {
	ctx *Context
}

func NewStepDiscoverPackage(ctx *Context) *StepDiscoverPackage {
	return &StepDiscoverPackage{ctx: ctx}
}

func (s *StepDiscoverPackage) Process(modulePath string) {
	files := fs.DiscoverModules(modulePath)
	for _, modulePath := range files {
		if !s.ctx.PreRegisterModule(modulePath) {
			continue
		}
		s.ctx.ToPrepareAST <- modulePath
	}
}
