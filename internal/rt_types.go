package internal

import (
	"fmt"
	"strings"
)

var _type_id = uint64(0)
var Void, Bool, Int, Float, String RtType

func init() {
	Void = NewPrimitive("Void", func() AstData { panic("Void does not have a default value") })
	Bool = NewPrimitive("Bool", func() AstData { return &AstBool{Value: false} })
	Int = NewPrimitive("Int", func() AstData { return &AstInt{Value: 0} })
	Float = NewPrimitive("Float", func() AstData { return &AstFloat{Value: 0.0} })
	String = NewPrimitive("String", func() AstData { return &AstString{Value: ""} })
}

// Base Type
type baseType struct {
	id uint64
}

func newBaseType() baseType {
	_type_id++
	return baseType{_type_id}
}

func (t baseType) Id() uint64 { return t.id }

// Primitive Type
type PrimitiveType struct {
	baseType
	name string
	def  func() AstData
}

func NewPrimitive(name string, def func() AstData) *PrimitiveType {
	_type_id++
	return &PrimitiveType{
		baseType: baseType{_type_id},
		name:     name,
		def:      def,
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

func (t *PrimitiveType) Default() AstData {
	return t.def()
}

// Function Type
type FunctionType struct {
	baseType
	args []RtType
	ret  RtType
}

func NewFunctionType(args []RtType, ret RtType) *FunctionType {
	return &FunctionType{
		baseType: newBaseType(),
		args:     args,
		ret:      ret,
	}
}

func (t *FunctionType) Name() string {
	args := []string{}
	for _, arg := range t.args {
		args = append(args, arg.Name())
	}
	return f("Fn(%s) %s", strings.Join(args, ", "), t.ret.Name())
}

func (t *FunctionType) Accepts(other RtType) bool {
	if t == nil || other == nil {
		return false
	}
	fn, ok := other.(*FunctionType)
	if !ok {
		return false
	}

	if len(t.args) != len(fn.args) {
		return false
	}

	for i, arg := range t.args {
		if !arg.Accepts(fn.args[i]) {
			return false
		}
	}

	return t.ret.Accepts(fn.ret)
}

func (t *FunctionType) Default() AstData {
	panic("Function does not have a default value")
}

func (t *FunctionType) Apply(args []RtType) (RtType, error) {
	if len(t.args) != len(args) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(t.args), len(args))
	}

	for i, arg := range t.args {
		if !arg.Accepts(args[i]) {
			return nil, fmt.Errorf("expected argument %d to be %s, got %s", i, arg.Name(), args[i].Name())
		}
	}

	return t.ret, nil

}

// Module Type
type ModuleType struct {
	baseType
	name string
}

func NewModuleType(name string) *ModuleType {
	return &ModuleType{
		baseType: newBaseType(),
		name:     name,
	}
}

func (t *ModuleType) Name() string              { return t.name }
func (t *ModuleType) Accepts(other RtType) bool { return false }
func (t *ModuleType) Default() AstData          { panic("Module does not have a default value") }

// // Data Type
// type DataType struct {
// 	baseType
// 	name string
// }

// func NewDataType(name string) *DataType {
// 	return &DataType{
// 		baseType: newBaseType(),
// 		name:     name,
// 	}
// }

// func (t *DataType) Name() string { return t.name }

// func (t *DataType) Accepts(other RtType) bool {
// 	if t == nil || other == nil {
// 		return false
// 	}
// 	dt, ok := other.(*DataType)
// 	if !ok {
// 		return false
// 	}

// 	return t.id == dt.id
// }

// func (t *DataType) Default() AstData {
// 	panic("Data does not have a default value")
// }
