package build

import (
	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/semantic/types"
)

type StepFinish struct {
	ctx *Context
}

func NewStepFinish(ctx *Context) *StepFinish {
	return &StepFinish{ctx: ctx}
}

func (s *StepFinish) Process() {
	s.checkMainFunction()
	s.ctx.Done <- nil
}

func (s *StepFinish) checkMainFunction() {
	main, _ := s.ctx.Modules.Get(s.ctx.EntryModulePath)

	mainFunc := main.Scope.GetValue("main")
	if mainFunc == nil {
		panic("function 'main' not found")
	}

	mainFuncType := mainFunc.Type().(*types.Function)
	if mainFuncType.Return != semantic.Void {
		panic("function 'main' must not return any value")
	}
	if len(mainFuncType.Parameters) > 0 {
		panic("function 'main' must not have any parameter")
	}
}
