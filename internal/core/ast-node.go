package core

import (
	"fmt"

	"github.com/renatopp/golden/lang"
)

var _node_id = uint64(0)

// Represents a generic node in the Abstract Syntax Tree (AST).
type AstNode struct {
	id     uint64
	module *Module     // Module that this node was defined.
	token  *lang.Token // Tokens attached to the node
	data   AstData     // Specific AST information for this node.
	type_  TypeData    // Type annotation that represents the return type of this node.
	ref    *Ref        // The node assignment reference name in the IR
}

func NewNode(token *lang.Token, data AstData) *AstNode {
	_node_id++
	return &AstNode{
		id:    _node_id,
		token: token,
		data:  data,
	}
}

func NewEmptyNode() *AstNode {
	_node_id++
	return &AstNode{
		id: _node_id,
	}
}

func (n *AstNode) Id() uint64         { return n.id }
func (n *AstNode) Package() *Package  { return n.module.Package }
func (n *AstNode) Module() *Module    { return n.module }
func (n *AstNode) Token() *lang.Token { return n.token }
func (n *AstNode) Data() AstData      { return n.data }
func (n *AstNode) Type() TypeData     { return n.type_ }
func (n *AstNode) Ref() *Ref          { return n.ref }

func (n *AstNode) WithModule(module *Module) *AstNode {
	n.module = module
	return n
}

func (n *AstNode) WithToken(token *lang.Token) *AstNode {
	n.token = token
	return n
}

func (n *AstNode) WithData(data AstData) *AstNode {
	n.data = data
	return n
}

func (n *AstNode) WithType(tp TypeData) *AstNode {
	n.type_ = tp
	return n
}

func (n *AstNode) WithRef(ref *Ref) *AstNode {
	n.ref = ref
	return n
}

func (n *AstNode) Copy() *AstNode {
	_node_id++
	return &AstNode{
		id:     _node_id,
		module: n.module,
		token:  n.token,
		data:   n.data,
		type_:  n.type_,
	}
}

func (n *AstNode) ExpressionKind() ExpressionKind {
	if n.data == nil {
		return InvalidExpression
	}
	return n.data.ExpressionKind()
}

func (n *AstNode) Tag() string {
	dt := "<nil>"
	if n.data != nil {
		dt = n.data.Tag()
	}
	tp := ""
	if n.type_ != nil {
		tp = " \u2192 " + n.type_.Tag()
	}
	return fmt.Sprintf("%s%s", dt, tp)
}

func (n *AstNode) Signature() string {
	dt := "<nil>"
	if n.data != nil {
		dt = n.data.Signature()
	}
	tp := ""
	if n.type_ != nil {
		tp = "::" + n.type_.Signature()
	}
	return fmt.Sprintf("%s%s", dt, tp)
}

func (n *AstNode) Children() []*AstNode {
	if n.data == nil {
		return []*AstNode{}
	}
	return n.data.Children()
}

func (n *AstNode) Traverse(fn func(*AstNode, int)) {
	n.traverse(fn, 0)
}

func (n *AstNode) traverse(fn func(*AstNode, int), depth int) {
	if n == nil {
		return
	}

	fn(n, depth)
	for _, child := range n.Children() {
		child.traverse(fn, depth+1)
	}
}
