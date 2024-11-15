package core

type IrComp interface {
	Node() *AstNode
	Tag() string
}

type IrWriter interface {
	EnterModule(*Module)
	ExitModule()

	NewInt(int64, *AstNode) IrComp
}
