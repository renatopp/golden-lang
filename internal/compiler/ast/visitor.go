package ast

type Visitor interface {
	VisitModule(*Module)
}
