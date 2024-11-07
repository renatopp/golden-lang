package internal

import "strings"

var (
	Void   = createPrimitive("Void")
	Bool   = createPrimitive("Bool")
	Int    = createPrimitive("Int")
	Float  = createPrimitive("Float")
	String = createPrimitive("String")
)

// Primitive
type PrimitiveType struct {
	name string
}

func NewPrimitiveType(name string) *PrimitiveType {
	return &PrimitiveType{name}
}

func (t *PrimitiveType) Name() string { return t.name }

// Function
type FunctionType struct {
	Params     []*Node
	ReturnType *Node
}

func (t *FunctionType) Name() string {
	r := "Fn("
	params := []string{}
	for _, p := range t.Params {
		params = append(params, p.Type.Name())
	}
	r += strings.Join(params, ", ")
	r += ")"
	if t.ReturnType != nil {
		r += " " + t.ReturnType.Type.Name()
	}
	return r
}

func createPrimitive(name string) *Node {
	return NewEmptyNode().WithType(NewPrimitiveType(name))
}
