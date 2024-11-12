package build

type StepFinish struct {
	ctx *Context
}

func NewStepFinish(ctx *Context) *StepFinish {
	return &StepFinish{ctx: ctx}
}

func (s *StepFinish) Process(modulePath string) {
	// process
}
