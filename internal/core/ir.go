package core

type IrComp interface {
	Node() *AstNode
	Tag() string
}

type IrWriter interface {
	EnterModule(*Module)
	ExitModule()

	// Declare(string, IrComp, *AstNode) IrComp
	Int(int64, *AstNode)
	Float(float64, *AstNode)
	Bool(bool, *AstNode)
	String(string, *AstNode)
}
