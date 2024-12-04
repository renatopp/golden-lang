package optimizations

import "github.com/renatopp/golden/internal/compiler/ast"

type Pipeline struct {
	Phases []ast.Visitor
}

func NewPipeline(phases ...ast.Visitor) *Pipeline {
	return (&Pipeline{
		Phases: make([]ast.Visitor, 0),
	}).WithPhases(phases...)
}

func (p *Pipeline) WithPhases(phases ...ast.Visitor) *Pipeline {
	p.Phases = append(p.Phases, phases...)
	return p
}

func (p *Pipeline) Run(node ast.Node) ast.Node {
	for _, phase := range p.Phases {
		node = node.Visit(phase)
	}
	return node
}
