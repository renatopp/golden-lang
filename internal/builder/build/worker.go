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
		case modulePath := <-ctx.toDiscoverPackage:
			steps.DiscoverPackage.Process(modulePath)
		case modulePath := <-ctx.toPrepareAST:
			steps.PrepareAst.Process(modulePath)
		}

		if ctx.CanProceedToDependencyGraph() {
			packages := steps.DependencyGraph.Process()
			steps.ResolveBindings.Process(packages)
			steps.Finish.Process()
		}
	}
}
