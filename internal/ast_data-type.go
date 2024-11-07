package internal

import (
	"strings"

	"github.com/renatopp/golden/lang"
)

// Data Declaration
type AstDataDecl struct {
	Name         string
	Constructors []*DataConstructor
}

func (a *AstDataDecl) Kind() string { return "type" }

func (a *AstDataDecl) String() string {
	r := f("data %s = ", a.Name)
	constr := []string{}
	for _, c := range a.Constructors {
		fields := []string{}
		for _, field := range c.Fields {
			fields = append(fields, f("%s %s", field.Name, ident(field.Type.String(), 2)))
		}
		constr = append(constr, f("%s(%s)", c.Name, strings.Join(fields, ", ")))
	}
	return r + strings.Join(constr, " | ")
}

type DataConstructor struct {
	Token  *lang.Token
	Name   string
	Shape  string // unit, tuple or record
	Fields []*DataConstructorField
}

type DataConstructorField struct {
	Token *lang.Token
	Name  string
	Type  *Node // Type Expression
}

// Type Identifier
type AstTypeIdent struct {
	Name string
}

func (a *AstTypeIdent) Kind() string   { return "type" }
func (a *AstTypeIdent) String() string { return f("%s", a.Name) }

// Function Type
type AstFunctionType struct {
	Params     []*FunctionTypeParam
	ReturnType *Node // Type Expression
}

func (a *AstFunctionType) Kind() string { return "type" }
func (a *AstFunctionType) String() string {
	params := []string{}
	for _, arg := range a.Params {
		params = append(params, f("%s", ident(arg.Type.String(), 1)))
	}
	return f("Fn (%s) %s", strings.Join(params, ", "), ident(a.ReturnType.String(), 1))
}

type FunctionTypeParam struct {
	Type *Node // Type Expression
}
