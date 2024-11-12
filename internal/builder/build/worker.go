package build

type steps struct {
	DiscoverPackage *StepDiscoverPackage
	PrepareAst      *StepPrepareAst
	DependencyGraph *StepDependencyGraph
	ResolveBindings *StepResolveBindings
	Finish          *StepFinish
}

func createSteps(ctx *Context) *steps {
	return &steps{
		DiscoverPackage: NewStepDiscoverPackage(ctx),
		PrepareAst:      NewStepPrepareAst(ctx),
		DependencyGraph: NewStepDependencyGraph(ctx),
		ResolveBindings: NewStepResolveBindings(ctx),
		Finish:          NewStepFinish(ctx),
	}
}

func startWorker(steps *steps, ctx *Context) {
	for {
		select {
		case modulePath := <-ctx.ToDiscoverPackage:
			steps.DiscoverPackage.Process(modulePath)
		case modulePath := <-ctx.ToPrepareAST:
			steps.PrepareAst.Process(modulePath)
		case modulePath := <-ctx.ToDependencyGraph:
			steps.DependencyGraph.Process(modulePath)
		case modulePath := <-ctx.ToResolveBindings:
			steps.ResolveBindings.Process(modulePath)
		case modulePath := <-ctx.ToFinish:
			steps.Finish.Process(modulePath)
		}
	}
}
