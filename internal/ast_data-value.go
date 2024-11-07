package internal

import (
	"strings"

	"github.com/renatopp/golden/lang"
)

// Variable Declaration
type AstVariableDecl struct {
	Name  string
	Type  *Node // nullable, type expression
	Value *Node // nullable, value expression
}

func (a *AstVariableDecl) Kind() string { return "value" }
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

// Function Declaration
type AstFunctionDecl struct {
	Name       string
	Params     []*FunctionParam
	ReturnType *Node // nullable, type expression
	Body       *Node // value expression
}

func (a *AstFunctionDecl) Kind() string { return "value" }
func (a *AstFunctionDecl) String() string {
	params := []string{}
	for _, p := range a.Params {
		params = append(params, f("%s %s", p.Name, ident(p.Type.String(), 1)))
	}

	type_ := ""
	if a.ReturnType != nil {
		type_ = ident(a.ReturnType.String(), 1)
	}

	return f("fn %s(%s) %s %s", a.Name, strings.Join(params, ", "), type_, ident(a.Body.String(), 1))
}

type FunctionParam struct {
	Name string
	Type *Node // nullable, type expression
}

// Block
type AstBlock struct {
	Expressions []*Node
}

func (a *AstBlock) Kind() string { return "value" }
func (a *AstBlock) String() string {
	expr := []string{}
	for _, n := range a.Expressions {
		expr = append(expr, ident(n.String(), 1))
	}
	return f("{ %s }", strings.Join(expr, "; "))
}

// Unary Operator
type AstUnaryOp struct {
	Operator string
	Right    *Node // value expression
}

func (a *AstUnaryOp) Kind() string { return "value" }
func (a *AstUnaryOp) String() string {
	return f("%s%s", a.Operator, ident(a.Right.String(), 1))
}

// Binary Operator
type AstBinaryOp struct {
	Operator string
	Left     *Node // value expression
	Right    *Node // value expression
}

func (a *AstBinaryOp) Kind() string { return "value" }
func (a *AstBinaryOp) String() string {
	return f("%s %s %s", ident(a.Left.String(), 1), a.Operator, ident(a.Right.String(), 1))
}

// Assignment
type AstAssignment struct {
	Operator string
	Left     *Node // value expression
	Right    *Node // value expression
}

func (a *AstAssignment) Kind() string { return "value" }
func (a *AstAssignment) String() string {
	return f("%s %s %s", ident(a.Left.String(), 1), a.Operator, ident(a.Right.String(), 1))
}

// Access
type AstAccess struct {
	Target   *Node // value expression
	Accessor string
}

func (a *AstAccess) Kind() string { return "value" }
func (a *AstAccess) String() string {
	return f("%s.%s", ident(a.Target.String(), 1), a.Accessor)
}

// Type Application
type AstAppl struct {
	Shape  string // unit, tuple, or record
	Target *Node  // nullable, type OR value expression
	Args   []*ApplArgument
}

func (a *AstAppl) Kind() string { return "value" }
func (a *AstAppl) String() string {
	args := []string{}
	for _, arg := range a.Args {
		if arg.Name != "" {
			args = append(args, f("%s=%s", arg.Name, ident(arg.Value.String(), 1)))
		} else {
			args = append(args, ident(arg.Value.String(), 1))
		}
	}

	return f("%s(%s)", ident(a.Target.String(), 1), strings.Join(args, ", "))
}

type ApplArgument struct {
	Token *lang.Token
	Name  string
	Value *Node // value expression
}

// Integer Literal
type AstInt struct {
	Value int64
}

func (a *AstInt) Kind() string   { return "value" }
func (a *AstInt) String() string { return f("%d", a.Value) }

// Float Literal
type AstFloat struct {
	Value float64
}

func (a *AstFloat) Kind() string   { return "value" }
func (a *AstFloat) String() string { return f("%f", a.Value) }

// Boolean Literal
type AstBool struct {
	Value bool
}

func (a *AstBool) Kind() string   { return "value" }
func (a *AstBool) String() string { return f("%t", a.Value) }

// String Literal
type AstString struct {
	Value string
}

func (a *AstString) Kind() string   { return "value" }
func (a *AstString) String() string { return f("%q", a.Value) }

// Variable Identifier
type AstVarIdent struct {
	Name string
}

func (a *AstVarIdent) Kind() string   { return "value" }
func (a *AstVarIdent) String() string { return a.Name }
