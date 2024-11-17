package comp

import (
	"fmt"

	"github.com/renatopp/golden/internal/core"
)

type ref interface {
	Name() string
}

type Base struct {
	node *core.AstNode
}

func NewBase(node *core.AstNode) *Base {
	return &Base{node: node}
}
func (c *Base) Node() *core.AstNode { return c.node }

// Int
type Int struct {
	Base
	Value int64
}

func (c *Int) Tag() string { return fmt.Sprintf("int:%d", c.Value) }

var _ core.IrComp = &Int{}

// Float
type Float struct {
	Base
	Value float64
}

func (c *Float) Tag() string { return fmt.Sprintf("float:%f", c.Value) }

var _ core.IrComp = &Float{}

// Bool
type Bool struct {
	Base
	Value bool
}

func (c *Bool) Tag() string { return fmt.Sprintf("bool:%t", c.Value) }

var _ core.IrComp = &Bool{}

// String
type String struct {
	Base
	Value string
}

func (c *String) Tag() string { return fmt.Sprintf("string:%s", c.Value) }

var _ core.IrComp = &String{}

// Identifier
// type Identifier struct {
// 	*Base
// 	RefName ref
// }

// func (c *Identifier) Tag() string { return fmt.Sprintf("identifier:%s", c.NameRef.Name()) }

// Declare
type Declare struct {
	Base
	NameRef ref
	NameUid string
	Value   core.IrComp
}

func (c *Declare) Tag() string { return fmt.Sprintf("declare:%s", c.NameUid) }
