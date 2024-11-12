package build

type StepDiscoverPackage struct {
	ctx *Context
}

func NewStepDiscoverPackage(ctx *Context) *StepDiscoverPackage {
	return &StepDiscoverPackage{ctx: ctx}
}

func (s *StepDiscoverPackage) Process(modulePath string) {
	// process
}
