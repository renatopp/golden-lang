package ast

type Visitor interface {
	VisitModule(*Module)
	VisitImport(*Import)
	VisitString(*String)
	VisitVarIdent(*VarIdent)
}
