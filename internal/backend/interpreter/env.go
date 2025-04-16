package interpreter

import "fmt"

type ObjectKind string

const ObjectType = ObjectKind("type")
const ObjectValue = ObjectKind("value")

type Object interface {
	Kind() ObjectKind // Kind of the object (type or value)
	Id() uint64       // Object identifier
	TypeId() uint64   // Id of the type object
	Type() *Object    // Type instance of the object
}

type Assignment struct {
	// Capability ...
	Object *Object
}

type Env struct {
	parent *Env
	types  map[string]Assignment
	values map[string]Assignment
}

func NewEnv() *Env {
	return &Env{
		parent: nil,
		types:  make(map[string]Assignment),
		values: make(map[string]Assignment),
	}
}

func (e *Env) Create() *Env {
	c := NewEnv()
	c.parent = e
	return c
}

func (e *Env) DeclareType(name string, obj *Object) error {
	if _, ok := e.types[name]; ok {
		return fmt.Errorf("type %s already declared", name)
	}
	e.types[name] = Assignment{Object: obj}
	return nil
}

func (e *Env) DeclareValue(name string, obj *Object) error {
	if _, ok := e.values[name]; ok {
		return fmt.Errorf("value %s already declared", name)
	}
	e.values[name] = Assignment{Object: obj}
	return nil
}
