package internal

import (
	"strings"

	"github.com/renatopp/golden/lang"
)

type AstModule struct {
	Imports   []*AstModuleImport
	Types     []*Node
	Functions []*Node
	Variables []*Node
}

type AstModuleImport struct {
	Path  string
	Alias string
}

func (a *AstModule) Kind() string      { return "value" }
func (a *AstModule) Label() string     { return f("module") }
func (a *AstModule) String() string    { return f("module") }
func (a *AstModule) Children() []*Node { return appendAll(a.Types, a.Functions, a.Variables) }

// Variable Declaration
type AstVariableDecl struct {
	Name  string
	Type  *Node // nullable, type expression
	Value *Node // nullable, value expression
}

func (a *AstVariableDecl) Kind() string  { return "value" }
func (a *AstVariableDecl) Label() string { return f("variable-decl %s", a.Name) }
func (a *AstVariableDecl) String() string {
	type_ := ""
	value_ := ""
	if a.Type != nil {
		type_ = " " + ident(a.Type.String(), 1)
	}
	if a.Value != nil {
		value_ = " = " + ident(a.Value.String(), 1)
	}

	return f("let %s%s%s", a.Name, type_, value_)
}
func (a *AstVariableDecl) Children() []*Node {
	children := []*Node{}
	if a.Type != nil {
		children = append(children, a.Type)
	}
	if a.Value != nil {
		children = append(children, a.Value)
	}
	return children
}

// Function Declaration
type AstFunctionDecl struct {
	Name       string
	Params     []*FunctionParam
	ReturnType *Node // nullable, type expression
	Body       *Node // value expression
}

func (a *AstFunctionDecl) Kind() string  { return "value" }
func (a *AstFunctionDecl) Label() string { return f("function-decl %s", a.Name) }
func (a *AstFunctionDecl) String() string {
	params := []string{}
	for _, p := range a.Params {
		params = append(params, f("%s %s", p.Name, ident(p.Type.String(), 1)))
	}

	type_ := ""
	if a.ReturnType != nil {
		type_ = " " + ident(a.ReturnType.String(), 1)
	}

	return f("fn %s(%s) %s %s", a.Name, strings.Join(params, ", "), type_, ident(a.Body.String(), 1))
}
func (a *AstFunctionDecl) Children() []*Node {
	children := []*Node{}
	for _, p := range a.Params {
		if p.Type != nil {
			children = append(children, p.Type)
		}
	}
	if a.ReturnType != nil {
		children = append(children, a.ReturnType)
	}
	children = append(children, a.Body)
	return children
}

type FunctionParam struct {
	Name string
	Type *Node // nullable, type expression
}

// Block
type AstBlock struct {
	Expressions []*Node
}

func (a *AstBlock) Kind() string  { return "value" }
func (a *AstBlock) Label() string { return f("block") }
func (a *AstBlock) String() string {
	expr := []string{}
	for _, n := range a.Expressions {
		expr = append(expr, ident(n.String(), 1))
	}
	return f("{ %s }", strings.Join(expr, "; "))
}
func (a *AstBlock) Children() []*Node { return a.Expressions }

// Unary Operator
type AstUnaryOp struct {
	Operator string
	Right    *Node // value expression
}

func (a *AstUnaryOp) Kind() string  { return "value" }
func (a *AstUnaryOp) Label() string { return f("unary-op %s", a.Operator) }
func (a *AstUnaryOp) String() string {
	return f("%s%s", a.Operator, ident(a.Right.String(), 1))
}
func (a *AstUnaryOp) Children() []*Node {
	return []*Node{a.Right}
}

// Binary Operator
type AstBinaryOp struct {
	Operator string
	Left     *Node // value expression
	Right    *Node // value expression
}

func (a *AstBinaryOp) Kind() string  { return "value" }
func (a *AstBinaryOp) Label() string { return f("binary-op %s", a.Operator) }
func (a *AstBinaryOp) String() string {
	return f("%s %s %s", ident(a.Left.String(), 1), a.Operator, ident(a.Right.String(), 1))
}
func (a *AstBinaryOp) Children() []*Node {
	return []*Node{a.Left, a.Right}
}

// Assignment
type AstAssignment struct {
	Operator string
	Left     *Node // value expression
	Right    *Node // value expression
}

func (a *AstAssignment) Kind() string  { return "value" }
func (a *AstAssignment) Label() string { return f("assignment %s", a.Operator) }
func (a *AstAssignment) String() string {
	return f("%s %s %s", ident(a.Left.String(), 1), a.Operator, ident(a.Right.String(), 1))
}
func (a *AstAssignment) Children() []*Node {
	return []*Node{a.Left, a.Right}
}

// Access
type AstAccess struct {
	Target   *Node // value expression
	Accessor string
}

func (a *AstAccess) Kind() string  { return "value" }
func (a *AstAccess) Label() string { return f("access %s", a.Accessor) }
func (a *AstAccess) String() string {
	return f("%s.%s", ident(a.Target.String(), 1), a.Accessor)
}
func (a *AstAccess) Children() []*Node {
	return []*Node{a.Target}
}

// Type Application
type AstApply struct {
	Shape  string // unit, tuple, or record
	Target *Node  // nullable, type OR value expression
	Args   []*ApplyArgument
}

func (a *AstApply) Kind() string  { return "value" }
func (a *AstApply) Label() string { return f("apply %s", a.Shape) }
func (a *AstApply) String() string {
	args := []string{}
	for _, arg := range a.Args {
		if arg.Name != "" {
			args = append(args, f("%s=%s", arg.Name, ident(arg.Value.String(), 1)))
		} else {
			args = append(args, ident(arg.Value.String(), 1))
		}
	}

	target := ""
	if a.Target != nil {
		target = ident(a.Target.String(), 1) + " "
	}
	return f("%s(%s)", target, strings.Join(args, ", "))
}
func (a *AstApply) Children() []*Node {
	children := []*Node{}
	if a.Target != nil {
		children = append(children, a.Target)
	}
	for _, arg := range a.Args {
		children = append(children, arg.Value)
	}
	return children
}

type ApplyArgument struct {
	Token *lang.Token
	Name  string
	Value *Node // value expression
}

// Integer Literal
type AstInt struct {
	Value int64
}

func (a *AstInt) Kind() string      { return "value" }
func (a *AstInt) Label() string     { return f("int %d", a.Value) }
func (a *AstInt) String() string    { return f("%d", a.Value) }
func (a *AstInt) Children() []*Node { return []*Node{} }

// Float Literal
type AstFloat struct {
	Value float64
}

func (a *AstFloat) Kind() string      { return "value" }
func (a *AstFloat) Label() string     { return f("float %f", a.Value) }
func (a *AstFloat) String() string    { return f("%f", a.Value) }
func (a *AstFloat) Children() []*Node { return []*Node{} }

// Boolean Literal
type AstBool struct {
	Value bool
}

func (a *AstBool) Kind() string      { return "value" }
func (a *AstBool) Label() string     { return f("bool %t", a.Value) }
func (a *AstBool) String() string    { return f("%t", a.Value) }
func (a *AstBool) Children() []*Node { return []*Node{} }

// String Literal
type AstString struct {
	Value string
}

func (a *AstString) Kind() string      { return "value" }
func (a *AstString) Label() string     { return f("string %s", a.Value) }
func (a *AstString) String() string    { return f("%q", a.Value) }
func (a *AstString) Children() []*Node { return []*Node{} }

// Variable Identifier
type AstVarIdent struct {
	Name string
}

func (a *AstVarIdent) Kind() string      { return "value" }
func (a *AstVarIdent) Label() string     { return f("var-ident %s", a.Name) }
func (a *AstVarIdent) String() string    { return a.Name }
func (a *AstVarIdent) Children() []*Node { return []*Node{} }
