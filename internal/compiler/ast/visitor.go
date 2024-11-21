package ast

type Visitor interface {
	VisitModule(*Module)
	VisitImport(*Import)
	VisitInt(*Int)
	VisitFloat(*Float)
	VisitBool(*Bool)
	VisitString(*String)
	VisitVarIdent(*VarIdent)
	VisitVarDecl(*VarDecl)
	VisitBlock(*Block)
	VisitUnaryOp(*UnaryOp)
	VisitBinaryOp(*BinaryOp)
	VisitAccess(*Access)
	VisitTypeIdent(*TypeIdent)
	VisitFuncType(*FuncType)
	VisitFuncTypeParam(*FuncTypeParam)
	VisitFuncDecl(*FuncDecl)
	VisitFuncDeclParam(*FuncDeclParam)
	VisitAppl(*Appl)
	VisitApplArg(*ApplArg)
}
