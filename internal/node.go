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

// Import
type AstImport struct {
	Path  string
	Alias string
}

func (a *AstImport) String() string {
	if a.Alias != "" {
		return f("import %s as %s", a.Path, a.Alias)
	}
	return f("import %s", a.Path)
}
func (a *AstImport) Children() []*Node { return []*Node{} }

// Data Declaration
type AstDataDecl struct {
	Name         string
	Constructors []*Node
}

func (a *AstDataDecl) String() string    { return f("data decl %s", a.Name) }
func (a *AstDataDecl) Children() []*Node { return a.Constructors }

type AstConstructor struct {
	Name   string
	Shape  string // unit, tuple or record
	Fields []*Node
}

func (a *AstConstructor) String() string    { return f("constructor %s", a.Name) }
func (a *AstConstructor) Children() []*Node { return a.Fields }

type AstField struct {
	Name string
	Type *Node
}

func (a *AstField) String() string    { return f("field %s", a.Name) }
func (a *AstField) Children() []*Node { return []*Node{a.Type} }

// Function Declaration
type AstFunctionDecl struct {
	Name       string
	Parameters []*Node
	ReturnType *Node
	Body       *Node
}

func (a *AstFunctionDecl) String() string { return f("function decl %s", a.Name) }
func (a *AstFunctionDecl) Children() []*Node {
	return append(a.Parameters, []*Node{a.ReturnType, a.Body}...)
}

// Variable Declaration
type AstVariableDecl struct {
	Name       string
	Type       *Node
	Expression *Node
}

func (a *AstVariableDecl) String() string    { return f("variable decl %s", a.Name) }
func (a *AstVariableDecl) Children() []*Node { return []*Node{a.Type, a.Expression} }

// Parameter
type AstParameter struct {
	Name string
	Type *Node
}

func (a *AstParameter) String() string    { return f("parameter %s", a.Name) }
func (a *AstParameter) Children() []*Node { return []*Node{a.Type} }

// Type Ref
type AstTypeRef struct {
	Name string
}

func (a *AstTypeRef) String() string    { return f("typeref %s", a.Name) }
func (a *AstTypeRef) Children() []*Node { return []*Node{} }

type AstFnTypeRef struct {
	Parameters []*Node
	ReturnType *Node
}

func (a *AstFnTypeRef) String() string    { return "typeref Fn" }
func (a *AstFnTypeRef) Children() []*Node { return append(a.Parameters, a.ReturnType) }

// Block
type AstBlock struct {
	Expressions []*Node
}

func (a *AstBlock) String() string    { return "block" }
func (a *AstBlock) Children() []*Node { return a.Expressions }

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

// Var Identifier
type AstVarIdent struct {
	Name string
}

func (a *AstVarIdent) String() string    { return f("varident %s", a.Name) }
func (a *AstVarIdent) Children() []*Node { return []*Node{} }

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
