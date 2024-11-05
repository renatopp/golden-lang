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
	return f("<%s>", value)
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

// Variable Declaration
type AstVariableDecl struct {
	Name string
	Expr *Node
}
