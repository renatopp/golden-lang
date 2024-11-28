package golang

import (
	"os"
	"path"
	"strings"

	"github.com/renatopp/golden/internal/compiler/codegen/core"
	"github.com/renatopp/golden/internal/helpers/fs"
)

var _ core.Writer = &GoWriter{}

const main_go = `package main

import "golden/root"

func main() {
	root.Main()
}
`

const go_mod = `module golden

go 1.23.0

require ()
`

type GoWriter struct {
	targetDirectory string
	pkg             *Package
	scope           []core.ScopedInstruction
}

func NewGoWriter(targetDirectory string) *GoWriter {
	return &GoWriter{
		targetDirectory: targetDirectory,
		pkg:             nil,
		scope:           []core.ScopedInstruction{},
	}
}

func (w *GoWriter) Start() {
	root := path.Join(w.targetDirectory, "root")
	fs.GuaranteeDirectoryExists(root)
}

func (w *GoWriter) End() {
	fileName := path.Join(w.targetDirectory, "main.go")
	os.WriteFile(fileName, []byte(main_go), 0644)

	fileName = path.Join(w.targetDirectory, "go.mod")
	os.WriteFile(fileName, []byte(go_mod), 0644)
}

func (w *GoWriter) pop() core.ScopedInstruction {
	top := w.scope[len(w.scope)-1]
	w.scope = w.scope[:len(w.scope)-1]
	return top
}

func (w *GoWriter) push(i core.ScopedInstruction) {
	w.scope = append(w.scope, i)
}

func (w *GoWriter) Pop() {
	switch i := w.pop().(type) {
	case *Package:
		fileName := path.Join(i.PackagePath, "package.go")
		os.WriteFile(fileName, []byte(i.Print(0)), 0644)
	}
}

func (w *GoWriter) PushPackage(packagePath string, imports []string) {
	packageName := fs.PackagePath2PackageName(packagePath)

	targetPackageName := strings.ReplaceAll(packageName, "@", "root")
	targetPackagePath := path.Join(w.targetDirectory, targetPackageName)

	w.pkg = &Package{
		scoped:      newScoped(),
		PackageName: targetPackageName,
		PackagePath: targetPackagePath,
		Imports:     imports,
	}
	w.push(w.pkg)
	fs.GuaranteeDirectoryExists(targetPackagePath)
}

func (w *GoWriter) PushModule() {}

func (w *GoWriter) PushVarDecl() {}

func (w *GoWriter) PushFuncDecl() {}

func (w *GoWriter) PushBlock() {}

func (w *GoWriter) PushInt() {}

func (w *GoWriter) PushFloat() {}

func (w *GoWriter) PushBool() {}

func (w *GoWriter) PushString() {}

func (w *GoWriter) PushVarIdent() {}

func (w *GoWriter) PushTypeIdent() {}

func (w *GoWriter) PushAccess() {}

func (w *GoWriter) PushFuncAppl() {}

func (w *GoWriter) PushBinaryOp() {}

func (w *GoWriter) PushUnaryOp() {}
