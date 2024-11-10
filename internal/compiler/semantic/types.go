package semantic

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
)

var _type_id = uint64(0)
var Void, Bool, Int, Float, String core.TypeData

func init() {
	Void = &VoidType{}
	Bool = NewPrimitive("Bool", func() core.AstData { return &ast.Bool{Value: false} })
	Int = NewPrimitive("Int", func() core.AstData { return &ast.Int{Value: 0} })
	Float = NewPrimitive("Float", func() core.AstData { return &ast.Float{Value: 0.0} })
	String = NewPrimitive("String", func() core.AstData { return &ast.String{Value: ""} })
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

type VoidType struct {
	baseType
}

func (t *VoidType) Tag() string       { return "Void" }
func (t *VoidType) Signature() string { return "Void" }
func (t *VoidType) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}
	return true
}
func (t *VoidType) Default() core.AstData { panic("Void does not have a default value") }

// Primitive Type
type PrimitiveType struct {
	baseType
	name string
	def  func() core.AstData
}

func NewPrimitive(name string, def func() core.AstData) *PrimitiveType {
	_type_id++
	return &PrimitiveType{
		baseType: baseType{_type_id},
		name:     name,
		def:      def,
	}
}

func (t *PrimitiveType) Tag() string       { return t.name }
func (t *PrimitiveType) Signature() string { return t.name }
func (t *PrimitiveType) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}
	prim, ok := other.(*PrimitiveType)
	if !ok {
		return false
	}

	return t.id == prim.id
}
func (t *PrimitiveType) Default() core.AstData {
	return t.def()
}

// Function Type
type FunctionType struct {
	baseType
	Args []core.TypeData
	Ret  core.TypeData
}

func NewFunctionType(args []core.TypeData, ret core.TypeData) *FunctionType {
	return &FunctionType{
		baseType: newBaseType(),
		Args:     args,
		Ret:      ret,
	}
}

func (t *FunctionType) Tag() string { return "Fn" }
func (t *FunctionType) Signature() string {
	args := []string{}
	for _, arg := range t.Args {
		args = append(args, arg.Signature())
	}
	return fmt.Sprintf("Fn(%s) %s", strings.Join(args, ", "), t.Ret.Signature())
}
func (t *FunctionType) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}
	fn, ok := other.(*FunctionType)
	if !ok {
		return false
	}

	if len(t.Args) != len(fn.Args) {
		return false
	}

	for i, arg := range t.Args {
		if !arg.Accepts(fn.Args[i]) {
			return false
		}
	}

	return t.Ret.Accepts(fn.Ret)
}
func (t *FunctionType) Default() core.AstData {
	panic("Function does not have a default value")
}
func (t *FunctionType) Apply(args []core.TypeData) (core.TypeData, error) {
	if len(t.Args) != len(args) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(t.Args), len(args))
	}

	for i, arg := range t.Args {
		if !arg.Accepts(args[i]) {
			return nil, fmt.Errorf("expected argument %d to be %s, got %s", i, arg.Signature(), args[i].Signature())
		}
	}

	return t.Ret, nil

}

// Module Type
type ModuleType struct {
	baseType
	Name   string
	Module *core.Module
}

func NewModuleType(name string, module *core.Module) *ModuleType {
	return &ModuleType{
		baseType: newBaseType(),
		Name:     name,
		Module:   module,
	}
}

func (t *ModuleType) Tag() string                      { return t.Name }
func (t *ModuleType) Signature() string                { return "" }
func (t *ModuleType) Accepts(other core.TypeData) bool { return false }
func (t *ModuleType) Default() core.AstData            { panic("Module does not have a default value") }

func (t *ModuleType) AccessValue(name string) (*core.AstNode, error) {
	val := t.Module.Scope.GetValue(name)
	if val == nil {
		return nil, fmt.Errorf("value %s not found", name)
	}
	if val.Type() == nil {
		t.Module.Analyzer.ResolveValue(val)
	}
	return val, nil
}
func (t *ModuleType) AccessType(name string) (core.TypeData, error) {
	val := t.Module.Scope.GetType(name)
	if val == nil {
		return nil, fmt.Errorf("type %s not found", name)
	}
	return val, nil
}

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

// func (t *DataType) Accepts(other core.TypeData) bool {
// 	if t == nil || other == nil {
// 		return false
// 	}
// 	dt, ok := other.(*DataType)
// 	if !ok {
// 		return false
// 	}

// 	return t.id == dt.id
// }

// func (t *DataType) Default() corecore.AstData {
// 	panic("Data does not have a default value")
// }
