package inst

type Instruction interface{}

type Ref struct {
	Package string // package path
	Module  string // module path
	Name    string
	SSA     int64
}

type VarDef struct {
	Name *Ref
	Expr Instruction
}

type Int struct {
	Value int64
}

type Float struct {
	Value float64
}

type Bool struct {
	Value bool
}

type String struct {
	Value string
}
