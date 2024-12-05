package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
)

type Feature uint64

const (
	Inherit Feature = iota
	Allowed
	Denied
)

type State struct {
	parent          *State
	node            ast.Node
	currentScope    *env.Scope
	currentModule   *ast.Module
	currentFunction *ast.FnDecl
	currentBlock    *ast.Block
	featureReturn   Feature
	featureBreak    Feature
	featureContinue Feature
}

func NewState() *State {
	return &State{
		parent:          nil,
		node:            nil,
		currentScope:    nil,
		currentModule:   nil,
		currentFunction: nil,
		currentBlock:    nil,
		featureReturn:   Denied,
		featureBreak:    Denied,
		featureContinue: Denied,
	}
}

func (s *State) New(node ast.Node) *State {
	return &State{
		parent:          s,
		node:            node,
		currentScope:    s.currentScope,
		currentModule:   s.currentModule,
		currentFunction: s.currentFunction,
		currentBlock:    s.currentBlock,
		featureReturn:   Inherit,
		featureBreak:    Inherit,
		featureContinue: Inherit,
	}
}

func (s *State) WithNode(node ast.Node) *State { s.node = node; return s }
func (s *State) Node() ast.Node                { return s.node }

func (s *State) WithScope(scope *env.Scope) *State { s.currentScope = scope; return s }
func (s *State) Scope() *env.Scope                 { return s.currentScope }

func (s *State) WithModule(module *ast.Module) *State { s.currentModule = module; return s }
func (s *State) Module() *ast.Module                  { return s.currentModule }

func (s *State) WithFunction(fn *ast.FnDecl) *State { s.currentFunction = fn; return s }
func (s *State) Function() *ast.FnDecl              { return s.currentFunction }

func (s *State) WithBlock(block *ast.Block) *State { s.currentBlock = block; return s }
func (s *State) Block() *ast.Block                 { return s.currentBlock }

func (s *State) EnableReturn() *State  { s.featureReturn = Allowed; return s }
func (s *State) DisableReturn() *State { s.featureReturn = Denied; return s }
func (s *State) CanReturn() bool {
	if s.featureReturn == Allowed {
		return true
	}
	if s.featureReturn == Denied {
		return false
	}
	if s.parent != nil {
		return s.parent.CanReturn()
	}
	return false
}

func (s *State) EnableBreak() *State  { s.featureBreak = Allowed; return s }
func (s *State) DisableBreak() *State { s.featureBreak = Denied; return s }
func (s *State) CanBreak() bool {
	if s.featureBreak == Allowed {
		return true
	}
	if s.featureBreak == Denied {
		return false
	}
	if s.parent != nil {
		return s.parent.CanBreak()
	}
	return false
}

func (s *State) EnableContinue() *State  { s.featureContinue = Allowed; return s }
func (s *State) DisableContinue() *State { s.featureContinue = Denied; return s }
func (s *State) CanContinue() bool {
	if s.featureContinue == Allowed {
		return true
	}
	if s.featureContinue == Denied {
		return false
	}
	if s.parent != nil {
		return s.parent.CanContinue()
	}
	return false
}
