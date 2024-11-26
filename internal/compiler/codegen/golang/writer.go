package golang

type Ref struct {
	Package    string // path
	Module     string // path
	Identifier string // var or type name
	Counter    uint64 // counter for unique names
}

type Writer struct {
	// file *file
	// function *function
}

func (w *Writer) Pop()           {}
func (w *Writer) PushFile()      {}
func (w *Writer) PushVarDecl()   {}
func (w *Writer) PushFnDecl()    {}
func (w *Writer) PushBlock()     {}
func (w *Writer) PushInt()       {}
func (w *Writer) PushFloat()     {}
func (w *Writer) PushBool()      {}
func (w *Writer) PushString()    {}
func (w *Writer) PushVarIdent()  {}
func (w *Writer) PushTypeIdent() {}
func (w *Writer) PushAccess()    {}
func (w *Writer) PushFuncAppl()  {}
func (w *Writer) PushBinaryOp()  {}
func (w *Writer) PushUnaryOp()   {}

type File struct {
}

func (f *File) Open()   {}
func (f *File) Close()  {}
func (f *File) String() {}
