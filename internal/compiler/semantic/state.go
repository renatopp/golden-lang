package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
)

type Feature uint64

type StateReturns struct {
	List []ast.Node
}

func NewStateReturns() *StateReturns {
	return &StateReturns{List: []ast.Node{}}
}

type State struct {
	parent          *State
	node            ast.Node
	currentModule   *ast.Module
	currentFunction *ast.FnDecl
	currentBlock    *ast.Block
	currentReturns  *StateReturns
}

func NewState() *State {
	return &State{
		parent:          nil,
		node:            nil,
		currentModule:   nil,
		currentFunction: nil,
		currentBlock:    nil,
		currentReturns:  NewStateReturns(),
	}
}

func (s *State) New(node ast.Node) *State {
	return &State{
		parent:          s,
		node:            node,
		currentModule:   s.currentModule,
		currentFunction: s.currentFunction,
		currentBlock:    s.currentBlock,
		currentReturns:  s.currentReturns,
	}
}

func (s *State) Node() ast.Node { return s.node }

func (s *State) WithModule(module *ast.Module) *State { s.currentModule = module; return s }
func (s *State) Module() *ast.Module                  { return s.currentModule }

func (s *State) WithFunction(fn *ast.FnDecl) *State {
	s.currentFunction = fn
	s.currentReturns = NewStateReturns()
	return s
}
func (s *State) Function() *ast.FnDecl { return s.currentFunction }

func (s *State) Returns() []ast.Node { return s.currentReturns.List }
func (s *State) AddReturn(node ast.Node) *State {
	s.currentReturns.List = append(s.currentReturns.List, node)
	return s
}

func (s *State) WithBlock(block *ast.Block) *State { s.currentBlock = block; return s }
func (s *State) Block() *ast.Block                 { return s.currentBlock }
