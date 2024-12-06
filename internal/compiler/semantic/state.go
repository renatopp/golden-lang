package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/errors"
)

type Feature uint64

type State struct {
	parent          *State
	node            ast.Node
	currentModule   *ast.Module
	currentFunction *ast.FnDecl
	currentBlock    *ast.Block
	currentReturns  []ast.Node
}

func NewState() *State {
	return &State{
		parent:          nil,
		node:            nil,
		currentModule:   nil,
		currentFunction: nil,
		currentBlock:    nil,
		currentReturns:  nil,
	}
}

func (s *State) New(node ast.Node) *State {
	s.parent.checkCircularInitialization(node)

	return &State{
		parent:          s,
		node:            node,
		currentModule:   s.currentModule,
		currentFunction: s.currentFunction,
		currentBlock:    s.currentBlock,
		currentReturns:  s.currentReturns,
	}
}

func (s *State) checkCircularInitialization(node ast.Node) {
	if s == nil || s.node == nil {
		return
	}
	if s.node.IsEqual(node) {
		errors.ThrowAtNode(node, errors.CircularReferenceError, "circular initialization detected")
	}
	if s.parent != nil {
		s.parent.checkCircularInitialization(node)
	}
}

func (s *State) Node() ast.Node { return s.node }

func (s *State) WithModule(module *ast.Module) *State { s.currentModule = module; return s }
func (s *State) Module() *ast.Module                  { return s.currentModule }

func (s *State) WithFunction(fn *ast.FnDecl) *State {
	s.currentFunction = fn
	s.currentReturns = []ast.Node{}
	return s
}
func (s *State) Function() *ast.FnDecl { return s.currentFunction }

func (s *State) Returns() []ast.Node { return s.currentReturns }
func (s *State) AddReturn(node ast.Node) *State {
	s.currentReturns = append(s.currentReturns, node)
	return s
}

func (s *State) WithBlock(block *ast.Block) *State { s.currentBlock = block; return s }
func (s *State) Block() *ast.Block                 { return s.currentBlock }
