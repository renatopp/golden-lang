package ast

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

var f = fmt.Sprintf

func appendAll[T any](arrays ...[]T) []T {
	var out []T
	for _, arr := range arrays {
		out = append(out, arr...)
	}
	return out
}

// Module
type Module struct {
	Imports   []*core.AstNode
	Types     []*core.AstNode
	Functions []*core.AstNode
	Variables []*core.AstNode
}

func (a *Module) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Module) Tag() string                         { return f("module") }
func (a *Module) Signature() string                   { return f("…") }
func (a *Module) Children() []*core.AstNode           { return appendAll(a.Types, a.Functions, a.Variables) }

var _ core.AstData = &Module{}

// Module Import
type ModuleImport struct {
	Path  string
	Alias string
}

func (a *ModuleImport) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *ModuleImport) Tag() string                         { return f("import:%s", a.Path) }
func (a *ModuleImport) Signature() string {
	if a.Alias != "" {
		return f("import '%s' as %s", a.Path, a.Alias)
	}
	return f("import '%s'", a.Path)
}
func (a *ModuleImport) Children() []*core.AstNode { return []*core.AstNode{} }

var _ core.AstData = &ModuleImport{}

// Variable Declaration
type VariableDecl struct {
	Name  string
	Type  *core.AstNode // nullable, type expression
	Value *core.AstNode // nullable, value expression
}

func (a *VariableDecl) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *VariableDecl) Tag() string                         { return f("variable-decl:%s", a.Name) }
func (a *VariableDecl) Signature() string {
	type_ := ""
	if a.Type != nil {
		type_ = f(" %s", a.Type.Signature())
	}
	value_ := ""
	if a.Value != nil {
		value_ = " …"
		// value_ = f(" = %s", a.Value.Signature())
	}
	return f("let %s%s%s", a.Name, type_, value_)
}
func (a *VariableDecl) Children() []*core.AstNode {
	children := []*core.AstNode{}
	if a.Type != nil {
		children = append(children, a.Type)
	}
	if a.Value != nil {
		children = append(children, a.Value)
	}
	return children
}

var _ core.AstData = &VariableDecl{}

// Function Declaration
type FunctionDecl struct {
	Name       string
	Params     []*FunctionDeclParam
	ReturnType *core.AstNode // nullable, type expression
	Body       *core.AstNode // value expression
}

type FunctionDeclParam struct {
	Name string
	Type *core.AstNode // nullable, type expression
}

func (a *FunctionDecl) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *FunctionDecl) Tag() string                         { return f("function-decl:%s", a.Name) }
func (a *FunctionDecl) Signature() string {
	params := []string{}
	for _, p := range a.Params {
		params = append(params, f("%s %s", p.Name, p.Type.Signature()))
	}
	ret := ""
	if a.ReturnType != nil {
		ret = f(" %s", a.ReturnType.Signature())
	}

	return f("fn %s(%s)%s %s", a.Name, strings.Join(params, ", "), ret, a.Body.Signature())
}
func (a *FunctionDecl) Children() []*core.AstNode {
	children := []*core.AstNode{}
	for _, p := range a.Params {
		if p.Type != nil {
			children = append(children, p.Type)
		}
	}
	if a.ReturnType != nil {
		children = append(children, a.ReturnType)
	}
	return append(children, a.Body)
}

var _ core.AstData = &FunctionDecl{}

// Block
type Block struct {
	Expressions []*core.AstNode
}

func (a *Block) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Block) Tag() string                         { return f("block") }
func (a *Block) Signature() string {
	if len(a.Expressions) == 0 {
		return "{}"
	}
	return "{ … }"
}
func (a *Block) Children() []*core.AstNode { return a.Expressions }

var _ core.AstData = &Block{}

// Unary Operator
type UnaryOp struct {
	Operator string
	Right    *core.AstNode // value expression
}

func (a *UnaryOp) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *UnaryOp) Tag() string                         { return f("unary-op:%s", a.Operator) }
func (a *UnaryOp) Signature() string                   { return f("%s…", a.Operator) }
func (a *UnaryOp) Children() []*core.AstNode           { return []*core.AstNode{a.Right} }

var _ core.AstData = &UnaryOp{}

// Binary Operator
type BinaryOp struct {
	Operator string
	Left     *core.AstNode // value expression
	Right    *core.AstNode // value expression
}

func (a *BinaryOp) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *BinaryOp) Tag() string                         { return f("binary-op:%s", a.Operator) }
func (a *BinaryOp) Signature() string                   { return f("… %s …", a.Operator) }
func (a *BinaryOp) Children() []*core.AstNode           { return []*core.AstNode{a.Left, a.Right} }

var _ core.AstData = &BinaryOp{}

// Access
type Access struct {
	Target   *core.AstNode // value expression
	Accessor string
}

func (a *Access) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Access) Tag() string                         { return f("access:%s", a.Accessor) }
func (a *Access) Signature() string                   { return f("….%s", a.Accessor) }
func (a *Access) Children() []*core.AstNode           { return []*core.AstNode{a.Target} }

var _ core.AstData = &Access{}

// Type Application
type Apply struct {
	Shape  string        // unit, tuple, or record
	Target *core.AstNode // nullable, type OR value expression
	Args   []*ApplyArgument
}

type ApplyArgument struct {
	Token *lang.Token
	Name  string
	Value *core.AstNode // value expression
}

func (a *Apply) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Apply) Tag() string                         { return f("apply:%s", a.Shape) }
func (a *Apply) Signature() string {
	args := []string{}
	for _, arg := range a.Args {
		if arg.Name != "" {
			args = append(args, f("%s=%s", arg.Name, arg.Value.Signature()))
		} else {
			args = append(args, arg.Value.Signature())
		}
	}

	target := ""
	if a.Target != nil {
		target = f("%s", a.Target.Signature())
	}
	return f("%s(%s)", target, strings.Join(args, ", "))
}
func (a *Apply) Children() []*core.AstNode {
	children := []*core.AstNode{}
	if a.Target != nil {
		children = append(children, a.Target)
	}
	for _, arg := range a.Args {
		children = append(children, arg.Value)
	}
	return children
}

var _ core.AstData = &Apply{}

// Integer Literal
type Int struct {
	Value int64
}

func (a *Int) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Int) Tag() string                         { return f("int:%d", a.Value) }
func (a *Int) Signature() string                   { return f("%d", a.Value) }
func (a *Int) Children() []*core.AstNode           { return []*core.AstNode{} }

// Float Literal
type Float struct {
	Value float64
}

func (a *Float) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Float) Tag() string                         { return f("float:%f", a.Value) }
func (a *Float) Signature() string                   { return f("%f", a.Value) }
func (a *Float) Children() []*core.AstNode           { return []*core.AstNode{} }

// Boolean Literal
type Bool struct {
	Value bool
}

func (a *Bool) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *Bool) Tag() string                         { return f("bool:%t", a.Value) }
func (a *Bool) Signature() string                   { return f("%t", a.Value) }
func (a *Bool) Children() []*core.AstNode           { return []*core.AstNode{} }

// String Literal
type String struct {
	Value string
}

func (a *String) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *String) Tag() string                         { return f("string:%s", a.Value) }
func (a *String) Signature() string                   { return f("%q", a.Value) }
func (a *String) Children() []*core.AstNode           { return []*core.AstNode{} }

// Variable Identifier
type VarIdent struct {
	Name string
}

func (a *VarIdent) ExpressionKind() core.ExpressionKind { return core.ValueExpression }
func (a *VarIdent) Tag() string                         { return f("var-ident:%s", a.Name) }
func (a *VarIdent) Signature() string                   { return a.Name }
func (a *VarIdent) Children() []*core.AstNode           { return []*core.AstNode{} }
