package build

type StepResolveBindings struct {
	ctx *Context
}

func NewStepResolveBindings(ctx *Context) *StepResolveBindings {
	return &StepResolveBindings{ctx: ctx}
}

func (s *StepResolveBindings) Process(modulePath string) {
	// process
}
