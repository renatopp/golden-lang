package internal

import "github.com/renatopp/golden/lang"

type AstData interface {
	String() string
	Children() []*Node
}

type Node struct {
	Token *lang.Token
	Data  AstData
}

func NewNode(token *lang.Token, data AstData) *Node {
	return &Node{Token: token, Data: data}
}

func (n *Node) ReplaceBy(node *Node) {
	n.Token = node.Token
	n.Data = node.Data
}
func (n *Node) WithToken(token *lang.Token) *Node {
	n.Token = token
	return n
}
func (n *Node) WithData(data AstData) *Node {
	n.Data = data
	return n
}
func (n *Node) String() string {
	if n == nil || n.Data == nil {
		return ""
	}
	value := n.Data.String()
	return f("[%s]", value)
}
func (n *Node) Children() []*Node {
	if n == nil || n.Data == nil {
		return []*Node{}
	}
	return n.Data.Children()
}
func (n *Node) Traverse(visitor func(*Node, int) bool) {
	n.traverse(visitor, 0)
}
func (n *Node) traverse(visitor func(*Node, int) bool, depth int) {
	if n == nil {
		return
	}
	if visitor(n, depth) {
		for _, child := range n.Children() {
			child.traverse(visitor, depth+1)
		}
	}
}

// Module
type AstModule struct {
	Imports   []*Node
	Types     []*Node
	Functions []*Node
	Variables []*Node
}

func (a *AstModule) String() string { return "module" }
func (a *AstModule) Children() []*Node {
	return appendAll(a.Imports, a.Types, a.Functions, a.Variables)
}

// Function Declaration
type AstFunctionDecl struct {
	Name string
	Body *Node
}

func (a *AstFunctionDecl) String() string    { return f("function %s", a.Name) }
func (a *AstFunctionDecl) Children() []*Node { return []*Node{a.Body} }

// Block
type AstBlock struct {
	Expressions []*Node
}

func (a *AstBlock) String() string    { return "block" }
func (a *AstBlock) Children() []*Node { return a.Expressions }

// Variable Declaration
type AstVariableDecl struct {
	Name string
	Expr *Node
}

// Int
type AstInt struct {
	Value int64
}

func (a *AstInt) String() string    { return f("int %d", a.Value) }
func (a *AstInt) Children() []*Node { return []*Node{} }

// Float
type AstFloat struct {
	Value float64
}

func (a *AstFloat) String() string    { return f("float %f", a.Value) }
func (a *AstFloat) Children() []*Node { return []*Node{} }

// String
type AstString struct {
	Value string
}

func (a *AstString) String() string    { return f("string %s", esc(a.Value)) }
func (a *AstString) Children() []*Node { return []*Node{} }

// Bool
type AstBool struct {
	Value bool
}

func (a *AstBool) String() string    { return f("bool %t", a.Value) }
func (a *AstBool) Children() []*Node { return []*Node{} }

// Unary
type AstUnary struct {
	Op    string
	Right *Node
}

func (a *AstUnary) String() string    { return f("unary %s", a.Op) }
func (a *AstUnary) Children() []*Node { return []*Node{a.Right} }

// Binary
type AstBinary struct {
	Op    string
	Left  *Node
	Right *Node
}

func (a *AstBinary) String() string    { return f("binary %s", a.Op) }
func (a *AstBinary) Children() []*Node { return []*Node{a.Left, a.Right} }
