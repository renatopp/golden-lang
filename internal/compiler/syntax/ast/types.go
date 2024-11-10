package ast

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

// Data Declaration
type DataDecl struct {
	Name         string
	Constructors []*DataConstructor
}

func (a *DataDecl) ExpressionKind() core.ExpressionKind { return core.TypeExpression }
func (a *DataDecl) Tag() string                         { return fmt.Sprintf("data-decl:%s", a.Name) }
func (a *DataDecl) Signature() string {
	// r := f("data %s = ", a.Name)
	// constr := []string{}
	// for _, c := range a.Constructors {
	// 	fields := []string{}
	// 	for _, field := range c.Fields {
	// 		fields = append(fields, f("%s %s", field.Name, ident(field.Type.String(), 2)))
	// 	}
	// 	constr = append(constr, f("%s(%s)", c.Name, strings.Join(fields, ", ")))
	// }
	// return r + strings.Join(constr, " | ")
	return ""
}
func (a *DataDecl) Children() []*core.AstNode {
	children := []*core.AstNode{}
	for _, c := range a.Constructors {
		for _, f := range c.Fields {
			children = append(children, f.Type)
		}
	}
	return children
}

var _ core.AstData = &DataDecl{}

type DataConstructor struct {
	Token  *lang.Token
	Name   string
	Shape  string // unit, tuple or record
	Fields []*DataConstructorField
}

type DataConstructorField struct {
	Token *lang.Token
	Name  string
	Type  *core.AstNode // Type Expression
}

// Type Identifier
var _ core.AstData = &TypeIdent{}

type TypeIdent struct {
	Name string
}

func (a *TypeIdent) ExpressionKind() core.ExpressionKind { return core.TypeExpression }
func (a *TypeIdent) Tag() string                         { return fmt.Sprintf("type-ident:%s", a.Name) }
func (a *TypeIdent) Signature() string                   { return fmt.Sprintf("%s", a.Name) }
func (a *TypeIdent) Children() []*core.AstNode           { return []*core.AstNode{} }

// Function Type
var _ core.AstData = &FunctionType{}

type FunctionType struct {
	Params     []*FunctionTypeParam
	ReturnType *core.AstNode // Nullable, type expression
}

type FunctionTypeParam struct {
	Type *core.AstNode // Type expression
}

func (a *FunctionType) ExpressionKind() core.ExpressionKind { return core.TypeExpression }
func (a *FunctionType) Tag() string                         { return fmt.Sprintf("function-type") }
func (a *FunctionType) Signature() string {
	params := []string{}
	for _, param := range a.Params {
		params = append(params, param.Type.Signature())
	}
	ret := ""
	if a.ReturnType != nil {
		ret = fmt.Sprintf(" %s", a.ReturnType.Signature())
	}
	return fmt.Sprintf("Fn(%s)%s", strings.Join(params, ", "), ret)
}
func (a *FunctionType) Children() []*core.AstNode {
	children := []*core.AstNode{}
	for _, arg := range a.Params {
		children = append(children, arg.Type)
	}
	children = append(children, a.ReturnType)
	return children
}
