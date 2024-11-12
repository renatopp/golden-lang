package build

type StepPrepareAst struct {
	ctx *Context
}

func NewStepPrepareAst(ctx *Context) *StepPrepareAst {
	return &StepPrepareAst{ctx: ctx}
}

func (s *StepPrepareAst) Process(modulePath string) {
	// process
}
