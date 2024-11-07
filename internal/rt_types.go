package internal

import "strings"

var _type_id = uint64(0)

var (
	Void   = createPrimitive("Void")
	Bool   = createPrimitive("Bool")
	Int    = createPrimitive("Int")
	Float  = createPrimitive("Float")
	String = createPrimitive("String")
)

// Primitive
type PrimitiveType struct {
	id   uint64
	name string
}

func NewPrimitiveType(name string) *PrimitiveType {
	_type_id++
	return &PrimitiveType{
		_type_id,
		name,
	}
}

func (t *PrimitiveType) Name() string { return t.name }

func (t *PrimitiveType) Accepts(other RtType) bool {
	if t == nil || other == nil {
		return false
	}
	prim, ok := other.(*PrimitiveType)
	if !ok {
		return false
	}

	return t.id == prim.id
}

func (t *PrimitiveType) Default() *Node {
	return nil
}

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

func (t *FunctionType) Accepts(other RtType) bool {
	return false
}

func (t *FunctionType) Default() *Node {
	return nil
}

func createPrimitive(name string) *Node {
	return NewEmptyNode().WithType(NewPrimitiveType(name))
}
