package core

type Instruction interface {
	Print(depth int) string
}

type ScopedInstruction interface {
	Instruction
	Append(Instruction)
}

type Writer interface {
	Start()
	End()
	Pop()
	PushPackage(path string, imports []string)
	PushModule()
	PushVarDecl()
	PushFuncDecl()
	PushBlock()
	PushInt()
	PushFloat()
	PushBool()
	PushString()
	PushVarIdent()
	PushTypeIdent()
	PushAccess()
	PushFuncAppl()
	PushBinaryOp()
	PushUnaryOp()
}
