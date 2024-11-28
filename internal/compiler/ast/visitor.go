package ast

type Visitor interface {
	VisitModule(Module) Node
	VisitConst(Const) Node
	VisitInt(Int) Node
	VisitFloat(Float) Node
	VisitString(String) Node
	VisitBool(Bool) Node
	VisitVarIdent(VarIdent) Node
	VisitTypeIdent(TypeIdent) Node
	VisitBinOp(BinOp) Node
	VisitUnaryOp(UnaryOp) Node
	VisitBlock(Block) Node
}
