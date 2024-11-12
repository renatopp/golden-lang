package build

type StepDependencyGraph struct {
	ctx *Context
}

func NewStepDependencyGraph(ctx *Context) *StepDependencyGraph {
	return &StepDependencyGraph{ctx: ctx}
}

func (s *StepDependencyGraph) Process(modulePath string) {
	// process
}
