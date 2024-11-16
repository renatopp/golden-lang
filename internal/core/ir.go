package core

type IrComp interface {
	Node() *AstNode
	Tag() string
}

type IrWriter interface {
	EnterModule(*Module)
	ExitModule()

	Declare(string, IrComp, *AstNode) IrComp
	NewInt(int64, *AstNode) IrComp
	NewFloat(float64, *AstNode) IrComp
	NewBool(bool, *AstNode) IrComp
	NewString(string, *AstNode) IrComp
}
